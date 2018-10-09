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
