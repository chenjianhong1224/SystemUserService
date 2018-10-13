package main

import (
	"fmt"

	"github.com/satori/go.uuid"
	"go.uber.org/zap"
)

type system_menu_service struct {
	d *dbOperator
}

func (m *system_menu_service) querySysMeunByExample(example SystemManagerMenuReqData) ([]TSysMenu, error) {
	args := []interface{}{}
	tmp := TSysMenu{}
	var sql string
	sql = "select Menu_id, Menu_uuid, Menu_name, Is_leaf, Parent_uuid, Menu_level, Link_path, Menu_ico, Sys_code, Open_type, Menu_status, Create_time, Create_user, Update_time, Update_user, Remark from t_sys_menu where 1=1 "
	if len(example.MenuId) != 0 {
		sql += " and Menu_id = ? "
		args = append(args, example.MenuId)
	}
	if example.MenuLevel != 0 {
		sql += " and Menu_level = ? "
		args = append(args, example.MenuLevel)
	}
	if len(example.MenuLink) != 0 {
		sql += " and Link_path = ? "
		args = append(args, example.MenuLink)
	}
	if len(example.MenuName) != 0 {
		sql += " and Menu_name = ? "
		args = append(args, example.MenuName)
	}
	if len(example.MenuParent) != 0 {
		sql += " and Parent_uuid = ? "
		args = append(args, example.MenuParent)
	}
	queryReq := &SqlQueryRequest{
		SQL:         sql,
		Args:        args,
		RowTemplate: tmp}
	reply := m.d.dbCli.Query(queryReq)
	queryRep, _ := reply.(*SqlQueryReply)
	if queryRep.Err != nil {
		zap.L().Error(fmt.Sprintf("query sys menu error:%s", queryRep.Err.Error()))
		return nil, queryRep.Err
	}
	var returnMenus []TSysMenu = []TSysMenu{}
	for i := 0; i < len(queryRep.Rows); i++ {
		returnMenus = append(returnMenus, queryRep.Rows[i].(TSysMenu))
	}
	return returnMenus, nil
}

func (m *system_menu_service) addSysMenu(systemManagerMenuReq SystemManagerMenuReqData, opUserId string) (string, error) {
	uid, _ := uuid.NewV4()
	args := []interface{}{}
	args = append(args, uid.String())
	args = append(args, systemManagerMenuReq.MenuName)
	args = append(args, systemManagerMenuReq.MenuParent)
	args = append(args, systemManagerMenuReq.MenuLevel)
	args = append(args, systemManagerMenuReq.MenuLink)
	args = append(args, opUserId)
	args = append(args, opUserId)
	execReq := SqlExecRequest{
		SQL:  "insert into t_sys_menu(Menu_uuid, Menu_name, Parent_uuid, Menu_level, Link_path, Menu_status, Create_time, Create_user, Update_time, Update_user) values (?, ?, ?, ?, ?, 1, now(), ?, now(), ?)",
		Args: args,
	}
	sqlReply := m.d.dbCli.Query(execReq)
	if sqlReply.Error() == nil {
		return uid.String(), nil
	}
	zap.L().Error(fmt.Sprintf("add sys menu[%s] error:%s", systemManagerMenuReq.MenuName, sqlReply.Error()))
	return "", sqlReply.Error()
}

func (m *system_menu_service) deleteSysMenu(menuId string, opUserId string) error {
	args := []interface{}{}
	args = append(args, opUserId)
	args = append(args, menuId)
	execReq := SqlExecRequest{
		SQL:  "update t_sys_menu set menu_status = 0, Update_time = now(), Update_User = ? where menu_uuid = ?",
		Args: args,
	}
	execReq2 := SqlExecRequest{
		SQL:  "update t_sys_role_menu set Status = 0, Update_time = now(), Update_User = ? where menu_uuid = ?",
		Args: args,
	}

	var execReqList = []SqlExecRequest{execReq, execReq2}
	err := m.d.dbCli.TransationExcute(execReqList)
	if err == nil {
		return nil
	}
	zap.L().Error(fmt.Sprintf("delete sys menu[%s] error:%s", menuId, err.Error()))
	return err
}

func (m *system_menu_service) updateSysMenu(systemManagerMenuReq SystemManagerMenuReqData, opUserId string) (string, error) {
	args := []interface{}{}
	args = append(args, systemManagerMenuReq.MenuName)
	args = append(args, systemManagerMenuReq.MenuParent)
	args = append(args, systemManagerMenuReq.MenuLevel)
	args = append(args, systemManagerMenuReq.MenuLink)
	args = append(args, opUserId)
	args = append(args, systemManagerMenuReq.MenuId)
	execReq := SqlExecRequest{
		SQL:  "update t_sys_menu set Menu_name = ?, Parent_uuid = ?, Menu_level = ?, Link_path = ?, Update_time = now(), Update_user = ? where menu_uuid = ?",
		Args: args,
	}
	sqlReply := m.d.dbCli.Query(execReq)
	if sqlReply.Error() == nil {
		return "", nil
	}
	reply := sqlReply.(*SqlExecReply)
	if reply.RowsAffected == 0 {
		return "未找到更新记录", nil
	}
	zap.L().Error(fmt.Sprintf("update sys menu[%s] error:%s", systemManagerMenuReq.MenuId, sqlReply.Error()))
	return "", sqlReply.Error()
}
