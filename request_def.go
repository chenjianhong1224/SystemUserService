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
	SysUserId    string                             `json:"sysUserId"`
	UserName     string                             `json:"userName"`
	LoginName    string                             `json:"loginName"`
	UserMobile   string                             `json:"userMobile"`
	UserEMail    string                             `json:"userEMail"`
	Passwd       string                             `json:"passwd"`
	HeadPortrait string                             `json:"headPortrait"`
	RoleList     []SystemManagerUserReqDataRoleList `json:"roleList"`
}

type SystemManagerUserReqDataRoleList struct {
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
