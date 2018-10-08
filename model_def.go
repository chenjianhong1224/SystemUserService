package main

import (
	"database/sql"

	"github.com/go-sql-driver/mysql"
)

type tWholeSaler struct {
	Saler_id     int64
	Saler_uuid   string
	Saler_name   sql.NullString
	Company      string
	Mobile       string
	Saler_status int32
	Create_time  mysql.NullTime
	Create_user  sql.NullString
	Update_time  mysql.NullTime
	Update_user  sql.NullString
	Remark       sql.NullString
}

type tUser struct {
	User_id       int64
	User_uuid     string
	User_name     sql.NullString
	Passwd        string
	Open_id       sql.NullString
	Other_from    sql.NullInt64
	Nickname      sql.NullString
	Head_portrait sql.NullString
	Agent_uuid    string
	User_type     int32
	User_status   int32
	User_token    sql.NullString
	Expiry_time   mysql.NullTime
	Create_time   mysql.NullTime
	Create_user   sql.NullString
	Update_time   mysql.NullTime
	Update_user   sql.NullString
	Remark        sql.NullString
}

type tWxPluginProgram struct {
	Program_id     string
	Program_uuid   string
	Program_name   string
	Appid          string
	Appsecrete     string
	Program_status int32
	Saler_uuid     sql.NullString
	Program_type   int32
}
