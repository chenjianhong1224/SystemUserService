package main

import (
	_ "github.com/go-sql-driver/mysql"
)

type dbOperator struct {
	dbCli *MysqlClient
}

func newDbOperator(cfg *Config) *dbOperator {
	dbConf := &DBConf{
		ConnectTimeoutSec:    cfg.Database.ConnTimeout,
		ReadTimeoutSec:       cfg.Database.ReadTimeout,
		WriteTimeoutSec:      cfg.Database.WriteTimeout,
		MaxOpenConns:         cfg.Database.MaxOpenConnNum,
		MaxIdleConns:         cfg.Database.MaxIdleConnNum,
		ConnKeepAliveTimeSec: cfg.Database.KeepAliveTime,
	}
	d := &dbOperator{}
	d.dbCli = NewMysqlClient(cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.DbName, dbConf)

	return d
}
