package main

import (
	"crypto/md5"
	"fmt"

	"github.com/satori/go.uuid"
	"go.uber.org/zap"
)

type system_user_service struct {
	d *dbOperator
}

func (m *system_user_service) addSysUser(systemManagerUser SystemManagerUserReqData, opUserId string) (string, error) {
	uid, _ := uuid.NewV4()
	args := []interface{}{}
	args = append(args, uid.String())
	args = append(args, systemManagerUser.UserName)
	args = append(args, systemManagerUser.LoginName)
	args = append(args, systemManagerUser.UserEMail)
	args = append(args, systemManagerUser.UserMobile)
	data := []byte(systemManagerUser.Passwd)
	has := md5.Sum(data)
	args = append(args, fmt.Sprintf("%x", has))
	args = append(args, systemManagerUser.HeadPortrait)
	args = append(args, opUserId)
	args = append(args, opUserId)
	execReq := SqlExecRequest{
		SQL:  "insert into t_sys_user(User_uuid, User_name, Login_name, User_email, User_phone, Login_passwd, Head_portrait, User_status, Create_time, Create_user, Update_time, Update_user) values (?, ?, ?, ?, ?, ?, ?, ?, 1, now(), ?, now(), ?)",
		Args: args,
	}
	var execReqList = []SqlExecRequest{execReq}
	for i := 0; i < len(systemManagerUser.RoleList); i++ {
		args2 := []interface{}{}
		args2 = append(args2, systemManagerUser.RoleList[i].RoleId)
		args2 = append(args2, uid.String())
		args2 = append(args2, opUserId)
		args2 = append(args2, opUserId)
		execReq2 := SqlExecRequest{
			SQL:  "insert into t_sys_role_user(Role_uuid, User_uuid, Status, Create_time, Create_user, Update_time, Update_user) values (?, ?, 1, now(), ?, now(), ?)",
			Args: args2,
		}
		execReqList = append(execReqList, execReq2)
	}
	err := m.d.dbCli.TransationExcute(execReqList)
	if err == nil {
		return uid.String(), nil
	}
	zap.L().Error(fmt.Sprintf("add sys user[%s] error:%s", systemManagerUser.UserName, err.Error()))
	return "", err
}

func (m *system_user_service) deleteSysUser(userId string, opUserId string) error {
	args := []interface{}{}
	args = append(args, opUserId)
	args = append(args, userId)
	execReq := SqlExecRequest{
		SQL:  "update t_sys_user set User_status = 0, Update_time = now(), Update_User = ? where User_uuid = ?",
		Args: args,
	}
	args2 := []interface{}{}
	args2 = append(args2, opUserId)
	args2 = append(args2, userId)
	execReq2 := SqlExecRequest{
		SQL:  "update t_sys_role_user set Status = 0, Update_time = now(), Update_User = ? where User_uuid = ?",
		Args: args2,
	}

	var execReqList = []SqlExecRequest{execReq, execReq2}
	err := m.d.dbCli.TransationExcute(execReqList)
	if err == nil {
		return nil
	}
	zap.L().Error(fmt.Sprintf("delete sys user[%s] error:%s", userId, err.Error()))
	return err
}

/**
该接口是全量更新接口，先更新用户表，然后把用户角色全失效，再按参数提供的角色进行生效更新或新增
**/
func (m *system_user_service) updateSysUser(systemManagerUser SystemManagerUserReqData, opUserId string) (string, error) {
	args := []interface{}{}
	args = append(args, systemManagerUser.SysUserId)
	tmp := tSysUser{}
	queryReq := &SqlQueryRequest{
		SQL:         "select 1 from t_sys_user where User_uuid = ?",
		Args:        args,
		RowTemplate: tmp}
	reply := m.d.dbCli.Query(queryReq)
	queryRep, _ := reply.(*SqlQueryReply)
	if queryRep.Err != nil {
		zap.L().Error(fmt.Sprintf("update sys user error:%s", queryRep.Err.Error()))
		return "数据库错误", queryRep.Err
	}
	if len(queryRep.Rows) == 0 {
		return "不存该系统用户" + systemManagerUser.SysUserId, nil
	}

	args3 := []interface{}{}
	args3 = append(args3, systemManagerUser.UserName)
	args3 = append(args3, systemManagerUser.LoginName)
	args3 = append(args3, systemManagerUser.UserEMail)
	args3 = append(args3, systemManagerUser.UserMobile)
	args3 = append(args3, systemManagerUser.LoginName)
	data := []byte(systemManagerUser.Passwd)
	has := md5.Sum(data)
	args3 = append(args3, fmt.Sprintf("%x", has))
	args3 = append(args3, systemManagerUser.HeadPortrait)
	args3 = append(args3, opUserId)
	args3 = append(args3, systemManagerUser.SysUserId)
	execReq3 := SqlExecRequest{
		SQL:  "update t_sys_user set User_name=?, Login_name=?, User_email=?, User_phone=?, Login_passwd=?, Head_portrait=?, Update_time=now(), Update_user=? where User_uuid = ?",
		Args: args3,
	}
	args4 := []interface{}{}
	args4 = append(args4, opUserId)
	args4 = append(args4, systemManagerUser.SysUserId)
	execReq4 := SqlExecRequest{
		SQL:  "update t_sys_role_user set Status = 0, Update_time = now(), Update_User = ? where User_uuid = ?",
		Args: args4,
	}
	var execReqList = []SqlExecRequest{execReq3, execReq4}
	for i := 0; i < len(systemManagerUser.RoleList); i++ {
		args6 := []interface{}{}
		args6 = append(args6, systemManagerUser.RoleList[i].RoleId)
		args6 = append(args6, systemManagerUser.SysUserId)
		tmp := tSysRoleUser{}
		queryReq := &SqlQueryRequest{
			SQL:         "select 1 from t_sys_role_user where Role_uuid = ? and User_uuid = ?",
			Args:        args6,
			RowTemplate: tmp}
		reply := m.d.dbCli.Query(queryReq)
		queryRep, _ := reply.(*SqlQueryReply)
		if queryRep.Err != nil {
			zap.L().Error(fmt.Sprintf("update sys user error:%s", queryRep.Err.Error()))
			return "数据库错误", queryRep.Err
		}
		if len(queryRep.Rows) != 0 {
			args5 := []interface{}{}
			args5 = append(args5, opUserId)
			args5 = append(args5, systemManagerUser.SysUserId)
			args5 = append(args5, systemManagerUser.RoleList[i].RoleId)
			execReq5 := SqlExecRequest{
				SQL:  "update t_sys_role_user set Status = 1, Update_time = now(), Update_User = ? where User_uuid = ? and Role_uuid = ?",
				Args: args5,
			}
			execReqList = append(execReqList, execReq5)
		} else {
			args5 := []interface{}{}
			args5 = append(args5, systemManagerUser.RoleList[i].RoleId)
			args5 = append(args5, systemManagerUser.SysUserId)
			args5 = append(args5, opUserId)
			args5 = append(args5, opUserId)
			execReq5 := SqlExecRequest{
				SQL:  "insert into t_sys_role_user(Role_uuid, User_uuid, Status, Create_user, Create_time, Update_time, Update_User) values (?, ?, 1, ?, now(), ?, now())",
				Args: args5,
			}
			execReqList = append(execReqList, execReq5)
		}
	}
	err := m.d.dbCli.TransationExcute(execReqList)
	if err != nil {
		zap.L().Error(fmt.Sprintf("update sys user error:%s", err.Error()))
		return "", err
	}
	return "", nil
}

func (m *system_user_service) querySysUser(loginName string) (*tSysUser, error) {
	args := []interface{}{}
	args = append(args, loginName)
	tmp := tSysUser{}
	queryReq := &SqlQueryRequest{
		SQL:         "select User_id, User_uuid, User_name, Login_name, User_email, User_phone, Login_passwd, Head_portrait, User_status, Create_time, Create_user, Update_time, Update_user from t_sys_user where Login_name = ?",
		Args:        args,
		RowTemplate: tmp}
	reply := m.d.dbCli.Query(queryReq)
	queryRep, _ := reply.(*SqlQueryReply)
	if queryRep.Err != nil {
		zap.L().Error(fmt.Sprintf("query sys user error:%s", queryRep.Err.Error()))
		return nil, queryRep.Err
	}
	if len(queryRep.Rows) == 0 {
		return nil, nil
	}
	return queryRep.Rows[0].(*tSysUser), nil
}
