package main

import (
	"fmt"

	"github.com/satori/go.uuid"
	"go.uber.org/zap"
)

type system_role_service struct {
	d *dbOperator
}

func (m *system_role_service) addSysRole(systemManagerRole SystemManagerRoleReqData, opUserId string) (string, error) {
	uid, _ := uuid.NewV4()
	args := []interface{}{}
	args = append(args, uid.String())
	args = append(args, systemManagerRole.RoleName)
	var isLeaf int32
	isLeaf = 0
	if len(systemManagerRole.RoleParent) != 0 {
		isLeaf = 1
	}
	args = append(args, isLeaf)
	args = append(args, systemManagerRole.RoleParent)
	args = append(args, systemManagerRole.RoleLevel)
	args = append(args, opUserId)
	args = append(args, opUserId)
	execReq := SqlExecRequest{
		SQL:  "insert into t_sys_role(Role_uuid, Role_name, Is_leaf, Parent_uuid, Role_level, Role_status, Create_time, Create_user, Update_time, Update_user) values (?, ?, ?, ?, ?, 1, now(), ?, now(), ?)",
		Args: args,
	}
	var execReqList = []SqlExecRequest{execReq}
	for i := 0; i < len(systemManagerRole.UserList); i++ {
		args2 := []interface{}{}
		args2 = append(args2, uid.String())
		args2 = append(args2, systemManagerRole.UserList[i].SysUserId)
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
	zap.L().Error(fmt.Sprintf("add sys role user[%s] error:%s", systemManagerRole.RoleName, err.Error()))
	return "", err
}

/**
当userList存在时则删除user和role的对应关系，当userList不存在时，则对整个role失效，并删除与其对应的所有关系
**/
func (m *system_role_service) deleteSysRole(systemManagerRole SystemManagerRoleReqData, opUserId string) error {
	var execReqList = []SqlExecRequest{}
	for i := 0; i < len(systemManagerRole.UserList); i++ {
		args2 := []interface{}{}
		args2 = append(args2, opUserId)
		args2 = append(args2, systemManagerRole.RoleId)
		args2 = append(args2, systemManagerRole.UserList[i].SysUserId)
		execReq2 := SqlExecRequest{
			SQL:  "update t_sys_role_user set Status=0, Update_time=now(), Update_user=? where role_uuid=? and user_uuid=?",
			Args: args2,
		}
		execReqList = append(execReqList, execReq2)
	}
	if len(systemManagerRole.UserList) == 0 {
		args := []interface{}{}
		args = append(args, opUserId)
		args = append(args, systemManagerRole.RoleId)
		execReq := SqlExecRequest{
			SQL:  "update t_sys_role set Role_status=0, Update_time=now(), Update_user=? where Role_uuid=?",
			Args: args,
		}
		execReqList = append(execReqList, execReq)
		args3 := []interface{}{}
		args3 = append(args3, opUserId)
		args3 = append(args3, systemManagerRole.RoleId)
		execReq3 := SqlExecRequest{
			SQL:  "update t_sys_role_user set Status=0, Update_time=now(), Update_user=? where role_uuid=?",
			Args: args3,
		}
		execReqList = append(execReqList, execReq3)
		execReq4 := SqlExecRequest{
			SQL:  "update t_sys_role_menu set Status=0, Update_time=now(), Update_user=? where role_uuid=?",
			Args: args3,
		}
		execReqList = append(execReqList, execReq4)
	}
	err := m.d.dbCli.TransationExcute(execReqList)
	return err
}

/**
该接口是全量更新接口，先更新用户表，然后把用户角色全失效，再按参数提供的角色进行生效更新或新增
**/
func (m *system_role_service) updateSysRole(systemManagerRole SystemManagerRoleReqData, opUserId string) error {
	var execReqList = []SqlExecRequest{}
	args := []interface{}{}
	args = append(args, opUserId)
	execReq := SqlExecRequest{
		SQL:  "update t_sys_role set Role_name=?, Is_leaf=?, Parent_uuid=?, Role_level=?, Update_time=now(), Update_user=? where Role_uuid=?",
		Args: args,
	}
	execReqList = append(execReqList, execReq)
	args1 := []interface{}{}
	args1 = append(args1, opUserId)
	execReq1 := SqlExecRequest{
		SQL:  "update t_sys_role_user set status=0, Update_time=now(), Update_user=? where Role_uuid=? and User_uuid=?",
		Args: args1,
	}
	execReqList = append(execReqList, execReq1)
	for i := 0; i < len(systemManagerRole.UserList); i++ {
		args2 := []interface{}{}
		args2 = append(args2, systemManagerRole.RoleId)
		args2 = append(args2, systemManagerRole.UserList[0].SysUserId)
		args2 = append(args2, opUserId)
		args2 = append(args2, opUserId)
		execReq2 := SqlExecRequest{
			SQL:  "insert into t_sys_role_user(Role_uuid, User_uuid, Status, Create_time, Create_user, Update_time, Update_user) values (?, ?, 1, now(), ?, now(), ?) ON DUPLICATE KEY UPDATE Status = 1",
			Args: args2,
		}
		execReqList = append(execReqList, execReq2)
	}
	err := m.d.dbCli.TransationExcute(execReqList)
	return err
}

func (m *system_role_service) querySysUser(roleName string) ([]tSysRole, error) {
	args := []interface{}{}
	args = append(args, roleName)
	tmp := tSysRole{}
	queryReq := &SqlQueryRequest{
		SQL:         "select Role_id, Role_uuid, Role_name, Is_leaf, Parent_uuid, Role_level, Role_status, Create_time, Create_user, Update_time, Update_user from t_sys_role where Role_name = ?",
		Args:        args,
		RowTemplate: tmp}
	reply := m.d.dbCli.Query(queryReq)
	queryRep, _ := reply.(*SqlQueryReply)
	if queryRep.Err != nil {
		zap.L().Error(fmt.Sprintf("query sys role error:%s", queryRep.Err.Error()))
		return nil, queryRep.Err
	}
	var returnRoles []tSysRole = []tSysRole{}
	for i := 0; i < len(queryRep.Rows); i++ {
		returnRoles = append(returnRoles, queryRep.Rows[i].(tSysRole))
	}
	return returnRoles, nil
}
