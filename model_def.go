package main

import (
	"database/sql"

	"github.com/go-sql-driver/mysql"
)

type TSysUser struct {
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
	Expire_time   mysql.NullTime
	Create_time   mysql.NullTime
	Create_user   sql.NullString
	Update_time   mysql.NullTime
	Update_user   sql.NullString
	Remark        sql.NullString
}

type TSysRole struct {
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

type TSysRoleMenu struct {
	Power_id    int64
	Power_uuid  string
	Role_uuid   string
	Menu_uuid   string
	Status      int32
	Create_time mysql.NullTime
	Create_user sql.NullString
	Update_time mysql.NullTime
	Update_user sql.NullString
}

type TSysRoleUser struct {
	Role_uuid   string
	User_uuid   string
	Status      int32
	Create_time mysql.NullTime
	Create_user sql.NullString
	Update_time mysql.NullTime
	Update_user sql.NullString
}

type TSysMenu struct {
	Menu_id     int64
	Menu_uuid   string
	Menu_name   string
	Is_leaf     int32
	Parent_uuid sql.NullString
	Menu_level  int32
	Link_path   sql.NullString
	Menu_ico    sql.NullString
	Sys_code    sql.NullString
	Open_type   sql.NullInt64
	Menu_status int32
	Create_time mysql.NullTime
	Create_user sql.NullString
	Update_time mysql.NullTime
	Update_user sql.NullString
	Remark      sql.NullString
}
