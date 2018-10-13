package main

import (
	"fmt"

	"github.com/satori/go.uuid"
	"go.uber.org/zap"
)

type system_privilege_service struct {
	d *dbOperator
}

func (m *system_privilege_service) addSysPrivilege(systemMangerPrivilegeData SystemManagerPrivilegeData, opUserId string) ([]SystemManagerPrivilegeRespData, error) {
	var execReqList = []SqlExecRequest{}
	var resp = []SystemManagerPrivilegeRespData{}
	for i := 0; i < len(systemMangerPrivilegeData.MenuList); i++ {
		uid, _ := uuid.NewV4()
		args := []interface{}{}
		args = append(args, uid.String())
		args = append(args, systemMangerPrivilegeData.RoleId)
		args = append(args, systemMangerPrivilegeData.MenuList[i])
		args = append(args, opUserId)
		args = append(args, opUserId)
		execReq := SqlExecRequest{
			SQL:  "insert into t_sys_role_menu(Power_uuid, Role_uuid, Menu_uuid, Status, Create_time, Create_user, Update_time, Update_user) values (?, ?, ?, 1, now(), ?, now(), ?)",
			Args: args,
		}
		execReqList = append(execReqList, execReq)
		var elem SystemManagerPrivilegeRespData
		elem = SystemManagerPrivilegeRespData{PowerId: uid.String(), RoleId: systemMangerPrivilegeData.RoleId, MenuId: systemMangerPrivilegeData.MenuList[i]}
		resp = append(resp, elem)
	}
	err := m.d.dbCli.TransationExcute(execReqList)
	if err == nil {
		return resp, nil
	}
	zap.L().Error(fmt.Sprintf("add sys privilege [%s] error:%s", systemMangerPrivilegeData.RoleId, err.Error()))
	return nil, err
}

func (m *system_privilege_service) deleteSysPrivilege(powerId string, opUserId string) error {
	args := []interface{}{}
	args = append(args, opUserId)
	args = append(args, powerId)
	execReq := SqlExecRequest{
		SQL:  "update t_sys_role_menu set Status = 0, Update_time = now(), Update_use =? where power_id = ?",
		Args: args,
	}
	sqlReply := m.d.dbCli.Query(execReq)
	if sqlReply.Error() == nil {
		return nil
	}
	zap.L().Error(fmt.Sprintf("delete sys privilege[%s] error:%s", powerId, sqlReply.Error()))
	return sqlReply.Error()
}

func (m *system_privilege_service) updateSysPrivilege(systemMangerPrivilegeData SystemManagerPrivilegeData, opUserId string) error {
	var execReqList = []SqlExecRequest{}
	if len(systemMangerPrivilegeData.MenuList) > 0 {
		args := []interface{}{}
		args = append(args, opUserId)
		args = append(args, systemMangerPrivilegeData.RoleId)
		execReq := SqlExecRequest{
			SQL:  "update t_sys_role_menu set Status = 0, Update_time = now(), Update_use =? where role_id = ?",
			Args: args,
		}
		execReqList = append(execReqList, execReq)
		for i := 0; i < len(systemMangerPrivilegeData.MenuList); i++ {
			uid, _ := uuid.NewV4()
			args2 := []interface{}{}
			args2 = append(args2, uid.String())
			args2 = append(args2, systemMangerPrivilegeData.RoleId)
			args2 = append(args2, systemMangerPrivilegeData.MenuList[i])
			args2 = append(args2, opUserId)
			args2 = append(args2, opUserId)
			args2 = append(args2, opUserId)
			execReq2 := SqlExecRequest{
				SQL:  "insert into t_sys_role_menu(Power_uuid, Role_uuid, Menu_uuid, Status, Create_time, Create_user, Update_time, Update_user) values (?, ?, ?, 1, now(), ?, now(), ?) ON DUPLICATE KEY UPDATE Status = 1, Update_time = now(), Update_user = ?",
				Args: args2,
			}
			execReqList = append(execReqList, execReq2)
		}
		err := m.d.dbCli.TransationExcute(execReqList)
		return err
	}
	return nil
}

func (m *system_privilege_service) querySysPrivilegeByExample(systemMangerPrivilegeData SystemManagerPrivilegeData) ([]TSysRoleMenu, error) {
	args := []interface{}{}
	var sql string
	sql = "select Power_id, Power_uuid, Role_uuid, Menu_uuid, Status, Create_time, Create_user, Update_time, Update_user from t_sys_role_menu where 1=1 "
	if len(systemMangerPrivilegeData.PowerId) != 0 {
		args = append(args, systemMangerPrivilegeData.PowerId)
		sql += " Power_uuid = ? "
	}
	if len(systemMangerPrivilegeData.RoleId) != 0 {
		args = append(args, systemMangerPrivilegeData.RoleId)
		sql += " Role_uuid = ? "
	}
	tmp := TSysRoleMenu{}
	queryReq := &SqlQueryRequest{
		SQL:         sql,
		Args:        args,
		RowTemplate: tmp}
	reply := m.d.dbCli.Query(queryReq)
	queryRep, _ := reply.(*SqlQueryReply)
	if queryRep.Err != nil {
		zap.L().Error(fmt.Sprintf("query sys role menu error:%s", queryRep.Err.Error()))
		return nil, queryRep.Err
	}
	var returnRoleMenus []TSysRoleMenu = []TSysRoleMenu{}
	for i := 0; i < len(queryRep.Rows); i++ {
		returnRoleMenus = append(returnRoleMenus, queryRep.Rows[i].(TSysRoleMenu))
	}
	return returnRoleMenus, nil
}
