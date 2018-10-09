package main

import (
	"database/sql"

	"github.com/go-sql-driver/mysql"
)

type tSysUser struct {
	User_id       int64
	User_uuid     string
	User_name     sql.NullString
	Login_name    sql.NullString
	User_email    sql.NullString
	User_phone    sql.NullString
	Login_passwd  string
	Head_portrait sql.NullString
	User_status   int32
	Login_time    mysql.NullTime
	Login_from    sql.NullString
	Expiry_time   mysql.NullTime
	Create_time   mysql.NullTime
	Create_user   sql.NullString
	Update_time   mysql.NullTime
	Update_user   sql.NullString
	Remark        sql.NullString
}

type tSysRole struct {
	Role_id     int64
	Role_uuid   string
	Role_name   sql.NullString
	Is_leaf     int32
	Parent_uuid sql.NullString
	Role_level  int32
	Role_status int32
	Create_time mysql.NullTime
	Create_user sql.NullString
	Update_time mysql.NullTime
	Update_user sql.NullString
	Remark      sql.NullString
}

type tSysRoleMenu struct {
	Role_uuid   string
	Menu_uuid   string
	Status      int32
	Create_time mysql.NullTime
	Create_user sql.NullString
	Update_time mysql.NullTime
	Update_user sql.NullString
}

type tSysRoleUser struct {
	Role_uuid   string
	User_uuid   string
	Status      int32
	Create_time mysql.NullTime
	Create_user sql.NullString
	Update_time mysql.NullTime
	Update_user sql.NullString
}
