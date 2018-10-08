package main

import (
	"fmt"
	"reflect"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type MysqlClient struct {
	db *sqlx.DB
}

type DBConf struct {
	ConnectTimeoutSec    int
	ReadTimeoutSec       int
	WriteTimeoutSec      int
	MaxOpenConns         int
	MaxIdleConns         int
	ConnKeepAliveTimeSec time.Duration
}

func NewMysqlClient(host string, port int, user, passwd, dbname string, conf *DBConf) *MysqlClient {
	db, err := sqlx.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?timeout=%ds&charset=utf8&readTimeout=%ds&writeTimeout=%ds&parseTime=true", user, passwd, host, port, dbname, conf.ConnectTimeoutSec, conf.ReadTimeoutSec, conf.WriteTimeoutSec))
	if err != nil {
		panic(err)
	}

	db.SetMaxOpenConns(conf.MaxOpenConns)
	db.SetMaxIdleConns(conf.MaxIdleConns)
	db.SetConnMaxLifetime(conf.ConnKeepAliveTimeSec)

	if err := db.Ping(); err != nil {
		panic(err)
	}

	client := &MysqlClient{
		db: db,
	}

	return client
}

type SqlRequest interface {
}

type SqlReply interface {
	Error() error
}

type SqlExecRequest struct {
	SQL  string
	Args []interface{}
}

type SqlExecReply struct {
	RowsAffected int64
	LastInsertId int64
	Err          error
}

func (ser *SqlExecReply) Error() error {
	return ser.Err
}

type SqlQueryRequest struct {
	SQL         string
	Args        []interface{}
	RowTemplate interface{}
}

type SqlQueryReply struct {
	Rows []interface{}
	Err  error
}

func (sqr *SqlQueryReply) Error() error {
	return sqr.Err
}

func (c *MysqlClient) TransationExcute(reqList []SqlExecRequest) error {
	zap.L().Debug(fmt.Sprintf("%+v", reqList))
	conn, err := c.db.Begin()
	if err != nil {
		return err
	}
	for _, req := range reqList {
		stm, err := conn.Prepare(req.SQL)
		if err != nil {
			conn.Rollback()
			return err
		}
		_, err = stm.Exec(req.Args...)
		if err != nil {
			conn.Rollback()
			return err
		}
	}
	conn.Commit()
	return nil
}

func (c *MysqlClient) Query(req SqlRequest) SqlReply {
	switch t := req.(type) {
	case *SqlExecRequest:
		zap.L().Debug(fmt.Sprintf("%+v", t))
		conn, err := c.db.Begin()
		result, err := c.db.Exec(t.SQL, t.Args...)
		if err != nil {
			conn.Rollback()
			return &SqlExecReply{Err: err}
		}
		var count, newID int64
		count, err = result.RowsAffected()
		if err != nil {
			conn.Rollback()
			return &SqlExecReply{count, newID, err}
		}
		newID, err = result.LastInsertId()
		if err != nil {
			conn.Rollback()
			return &SqlExecReply{count, newID, err}
		}
		conn.Commit()
		return &SqlExecReply{count, newID, err}

	case *SqlQueryRequest:
		zap.L().Debug(fmt.Sprintf("%+v", t))
		rows, err := c.db.Query(t.SQL, t.Args...)
		defer func() {
			if rows != nil {
				rows.Close()
			}
		}()

		if err != nil {
			return &SqlQueryReply{Err: err}
		}

		reply := &SqlQueryReply{Rows: []interface{}{}}
		for rows.Next() {
			args := []interface{}{}
			row := reflect.New(reflect.ValueOf(t.RowTemplate).Type())
			elem := row.Elem()
			for i := 0; i < elem.NumField(); i++ {
				args = append(args, elem.Field(i).Addr().Interface())
			}
			err := rows.Scan(args...)
			if err != nil {
				reply.Err = fmt.Errorf("Scan row error: %s", err)
				break
			}
			reply.Rows = append(reply.Rows, row.Interface())
		}
		return reply
	default:
		panic(fmt.Sprintf("unknown command: %v", t))
	}
}
