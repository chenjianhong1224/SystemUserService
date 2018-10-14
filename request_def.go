package main

type RequestHead struct {
	RequestId string `json:"requestId"`
	UserId    string `json:"userId"`
	Cmd       int32  `json:"cmd"`
}

type ResponseHead struct {
	RequestId string `json:"requestId"`
	ErrorCode int32  `json:"errorCode"`
	ErrorMsg  string `json:"errorMsg"`
	Cmd       int32  `json:"cmd"`
}

type SystemManagerUserReqData struct {
	SysUserId    string   `json:"sysUserId"`
	UserName     string   `json:"userName"`
	LoginName    string   `json:"loginName"`
	UserMobile   string   `json:"userMobile"`
	UserEMail    string   `json:"userEMail"`
	Passwd       string   `json:"passwd"`
	HeadPortrait string   `json:"headPortrait"`
	RoleList     []string `json:"roleList"`
}

type SystemManagerUserReq struct {
	RequestHead
	Data SystemManagerUserReqData `json:"data"`
}

type SystemManagerUserRespData struct {
	SysUserId  string `json:"sysUserId"`
	UserName   string `json:"userName"`
	LoginName  string `json:"loginName"`
	UserMobile string `json:"userMobile"`
	UserEMail  string `json:"userEMail"`
}

type SystemManagerUserResp struct {
	ResponseHead
	Data SystemManagerUserRespData `json:"data"`
}

type SystemManagerRoleReq struct {
	RequestHead
	Data SystemManagerRoleReqData `json:"data"`
}

type SystemManagerRoleReqData struct {
	RoleId     string   `json:"roleId"`
	RoleName   string   `json:"roleName"`
	RoleLevel  int32    `json:"roleLevel"`
	RoleParent string   `json:"roleParent"`
	UserList   []string `json:"userList"`
}

type SystemManagerRoleResp struct {
	ResponseHead
	Data SystemManagerRoleRespData `json:"data"`
}

type SystemManagerRoleRespData struct {
	RoleId    string `json:"roleId"`
	RoleName  string `json:"roleName"`
	RoleLevel int32  `json:"roleLevel"`
}

type SystemManagerMenuReq struct {
	RequestHead
	Data SystemManagerMenuReqData `json:"data"`
}

type SystemManagerMenuReqData struct {
	MenuId     string `json:"menuId"`
	MenuName   string `json:"menuName"`
	MenuLevel  int32  `json:"menuLevel"`
	MenuParent string `json:"menuParent"`
	MenuLink   string `json:"menuLink"`
}

type SystemManagerMenuResp struct {
	ResponseHead
	Data SystemManagerMenuRespData `json:"data"`
}

type SystemManagerMenuRespData struct {
	MenuId    string `json:"menuId"`
	MenuName  string `json:"menuName"`
	MenuLevel int32  `json:"menuLevel"`
}

type SystemManagerPrivilegeReq struct {
	RequestHead
	Data SystemManagerPrivilegeData `json:"data"`
}

type SystemManagerPrivilegeData struct {
	PowerId  string   `json:"powerId"`
	RoleId   string   `json:"roleId"`
	MenuList []string `json:"menuList"`
}

type SystemManagerPrivilegeRespData struct {
	PowerId string `json:"powerId"`
	RoleId  string `json:"roleId"`
	MenuId  string `json:"menuId"`
}

type SystemManagerPrivilegeResp struct {
	ResponseHead
	Data []SystemManagerPrivilegeRespData `json:"data"`
}
