package main

type RequestHead struct {
	RequestId string `json:"requestId"`
	UserId    int32  `json:"userId"`
	Cmd       int32  `json:"cmd"`
	WsId      string `json:"wsId"`
}

type ResponseHead struct {
	RequestId string `json:"requestId"`
	ErrorCode int32  `json:"errorCode"`
	ErrorMsg  string `json:"errorMsg"`
	Cmd       int32  `json:"cmd"`
}

type SystemManagerUserReqData struct {
	SysUserId    string                         `json:"sysUserId"`
	UserName     string                         `json:"userName"`
	LoginName    string                         `json:"loginName"`
	UserMobile   string                         `json:"userMobile"`
	UserEMail    string                         `json:"userEMail"`
	Passwd       string                         `json:"passwd"`
	HeadPortrait string                         `json:"headPortrait"`
	RoleList     []SystemManagerUserReqDataRole `json:"roleList"`
}

type SystemManagerUserReqDataRole struct {
	RoleId string `json:"roleId"`
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
	RoleId     string                           `json:"roleId"`
	RoleName   string                           `json:"roleName"`
	RoleLevel  int32                            `json:"roleLevel"`
	RoleParent string                           `json:"roleParent"`
	UserList   []SystemManagerRoleReqDataUserId `json:"userList"`
}

type SystemManagerRoleReqDataUserId struct {
	SysUserId string `json:"sysUserId"`
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
