package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"

	"go.uber.org/zap"
)

type clientInfo struct {
	ipStr string
	ipNum int32
	port  int32
}

type httpHandler struct {
	cfg          *Config
	systemUserSv *system_user_service
}

func (ci *clientInfo) inetAton() {
	ip := net.ParseIP(ci.ipStr)
	ci.ipNum = int32(binary.BigEndian.Uint32(ip.To4()))
}

func (m *httpHandler) start() error {
	//start http server
	s := &http.Server{
		Addr:           m.cfg.Server.Endpoint,
		Handler:        nil,
		ReadTimeout:    m.cfg.Server.HttpReadTimeout,
		WriteTimeout:   m.cfg.Server.HttpWriteTimeout,
		MaxHeaderBytes: int(m.cfg.Server.MaxHeadSize),
	}
	http.HandleFunc("/api", m.process)
	go s.ListenAndServe()

	return nil
}

func (m *httpHandler) ivalidResp(w http.ResponseWriter) {
	http.Error(w, http.StatusText(http.StatusInternalServerError),
		http.StatusInternalServerError)
}

func (m *httpHandler) getClientInfo(r *http.Request) *clientInfo {
	cliIp, cliPort, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		zap.L().Warn(fmt.Sprintf("userip: %q is not IP:port", r.RemoteAddr))
		return &clientInfo{ipNum: 0, port: 0}
	} else {
		zap.L().Debug(fmt.Sprintf("package from %s:%s", cliIp, cliPort))
		p, e := strconv.Atoi(cliPort)
		if e != nil {
			zap.L().Error(fmt.Sprintf("strconv Atoi port fail"))
			p = 0
		}

		ci := &clientInfo{
			ipStr: cliIp,
			port:  int32(p),
			ipNum: 0,
		}

		ci.inetAton()
		return ci
	}
}

func (m *httpHandler) process(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		zap.L().Info(fmt.Sprintf("get method not support, method:%s", r.Method))
		statObj.statHandler.StatCount(StatInvalidMethodReq)
		m.ivalidResp(w)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		statObj.statHandler.StatCount(StatReadBody)
		m.ivalidResp(w)
		return
	} else {
		zap.L().Debug(fmt.Sprintf("recv body len:%d content:%s", len(body), body))
		var req RequestHead
		err := json.Unmarshal(body, &req)
		if err != nil {
			zap.L().Error(fmt.Sprintf("json transfer error %s", err.Error()))
			m.ivalidResp(w)
			return
		}
		if req.Cmd == 1000 {
		} else if req.Cmd == 1002 {
		} else if req.Cmd == 1004 {
		} else if req.Cmd == 1006 {
		} else {
			var respHead ResponseHead
			respHead = ResponseHead{RequestId: req.RequestId, ErrorCode: 9999, Cmd: req.Cmd, ErrorMsg: "cmd不合法"}
			jsonData, err := json.Marshal(respHead)
			if err != nil {
				zap.L().Error(fmt.Sprintf("json transfer error %s", err.Error()))
				m.ivalidResp(w)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(jsonData))
			return
		}
	}
}

//func (m *httpHandler) wholeSalerRegister(body []byte, w http.ResponseWriter) {

//	var req wholeSalerRegisterReq
//	err := json.Unmarshal(body, &req)
//	if err != nil {
//		zap.L().Error(fmt.Sprintf("json transfer error %s", err.Error()))
//		m.ivalidResp(w)
//		return
//	}
//	var resp wholeSalerRegisterResp
//	tUser, err := m.usersv.queryUser(req.OpenId, req.UserId)
//	if err == nil {
//		if tUser == nil {
//			resp = wholeSalerRegisterResp{ResponseHead{RequestId: req.RequestId, ErrorCode: 1, ErrorMsg: "查不到对应的销售员userId=" + req.UserId + ",OpenId=" + req.OpenId, Cmd: 136}, wholeSalerRegisterRespData{}}
//		} else {
//			if tUser.User_type != 2 || tUser.User_status != 1 {
//				resp = wholeSalerRegisterResp{ResponseHead{RequestId: req.RequestId, ErrorCode: 1, ErrorMsg: "该销售员不存在或不是合法状态", Cmd: 136}, wholeSalerRegisterRespData{}}
//			} else {
//				tWholeSaler, err := m.wholesalersv.queryWholesaler(req.Data.WsMobile, req.Data.WsCompany)
//				if err == nil {
//					if tWholeSaler != nil {
//						resp = wholeSalerRegisterResp{ResponseHead{RequestId: req.RequestId, ErrorCode: 1, ErrorMsg: "该批发商已经注册", Cmd: 136}, wholeSalerRegisterRespData{}}
//					} else {
//						if req.Data.WsMobile == "" {
//							resp = wholeSalerRegisterResp{ResponseHead{RequestId: req.RequestId, ErrorCode: 1, ErrorMsg: "新增批发商手机号不能为空", Cmd: 136}, wholeSalerRegisterRespData{}}
//						} else {
//							uuid, passwd, err := m.wholesalersv.addWholesaler(req)
//							if err == nil {
//								resp = wholeSalerRegisterResp{ResponseHead{RequestId: req.RequestId, ErrorCode: 0, Cmd: 136}, wholeSalerRegisterRespData{WsId: uuid, WsName: req.Data.WsName, WsCompany: req.Data.WsCompany, WsMobile: req.Data.WsMobile, WsIdentityCode: passwd}}
//							} else {
//								resp = wholeSalerRegisterResp{ResponseHead{RequestId: req.RequestId, ErrorCode: 1, ErrorMsg: "新增批发商失败:" + err.Error(), Cmd: 136}, wholeSalerRegisterRespData{}}
//							}
//						}
//					}
//				}
//			}
//		}
//		data, err := json.Marshal(resp)
//		if err != nil {
//			zap.L().Error(fmt.Sprintf("json transfer error %s", err.Error()))
//			m.ivalidResp(w)
//			return
//		}
//		w.WriteHeader(http.StatusOK)
//		w.Write([]byte(data))
//		return
//	}
//	zap.L().Error(fmt.Sprintf("get saler_user error %s", err.Error()))
//	m.ivalidResp(w)
//	return
//}

//func (m *httpHandler) userLogin(body []byte, w http.ResponseWriter) {
//	var req userLoginReq
//	err := json.Unmarshal(body, &req)
//	if err != nil {
//		zap.L().Error(fmt.Sprintf("json transfer error %s", err.Error()))
//		m.ivalidResp(w)
//		return
//	}
//	errMsg, openId, _, err := m.getWxUserInfo(req.Data.SpId, req.Data.WxCode)
//	if err != nil {
//		zap.L().Error(fmt.Sprintf("userLogin error %s", err.Error()))
//		m.ivalidResp(w)
//		return
//	} else if errMsg != "" {
//		var respHead ResponseHead
//		respHead = ResponseHead{RequestId: req.RequestHead.RequestId, ErrorCode: 9999, ErrorMsg: errMsg, Cmd: 133}
//		jsonData, err := json.Marshal(respHead)
//		if err != nil {
//			zap.L().Error(fmt.Sprintf("json transfer error %s", err.Error()))
//			m.ivalidResp(w)
//			return
//		}
//		w.WriteHeader(http.StatusOK)
//		w.Write([]byte(jsonData))
//		return
//	}
//	var resp userLoginResp
//	if req.UserType != 1 { //零售商登录不需要密码，其他都需要
//		tUser, err := m.usersv.queryUserByPasswd(req.Data.Passwd, req.Data.LoginName, req.UserType)
//		if err == nil {
//			if tUser == nil {
//				resp = userLoginResp{ResponseHead{RequestId: req.RequestId, ErrorCode: 1, ErrorMsg: "用户" + req.Data.LoginName + "登录失败: 用户名或密码不正确", Cmd: 133}, userLoginRespData{}}
//			} else {
//				err = m.usersv.bindUser(openId, tUser.User_uuid)
//				// TODO 处理登录
//				resp = userLoginResp{ResponseHead{RequestId: req.RequestId, ErrorCode: 0, Cmd: 133}, userLoginRespData{OpenId: openId,
//					UserId:   tUser.User_uuid,
//					UserType: req.UserType,
//					UserName: tUser.User_name.String,
//					HeadIco:  ""}}

//			}
//		}
//	} else {
//		tUser, err := m.usersv.queryUserByOpenId(openId)
//		if err == nil {
//			var usrUUid string
//			if tUser == nil {
//				usrUUid, err = m.usersv.addRetailer(openId)
//				if err != nil {
//					zap.L().Error(fmt.Sprintf("login add addRetailer error %s", err.Error()))
//					m.ivalidResp(w)
//					return
//				}
//			} else {
//				// TODO 处理登录
//				usrUUid = tUser.User_uuid
//			}
//			resp = userLoginResp{ResponseHead{RequestId: req.RequestId, ErrorCode: 0, Cmd: 133}, userLoginRespData{OpenId: openId,
//				UserId:   usrUUid,
//				UserType: req.UserType,
//				UserName: "",
//				HeadIco:  ""}}
//		}
//	}
//	jsonData, err := json.Marshal(resp)
//	if err != nil {
//		zap.L().Error(fmt.Sprintf("json transfer error %s", err.Error()))
//		m.ivalidResp(w)
//		return
//	}
//	w.WriteHeader(http.StatusOK)
//	w.Write([]byte(jsonData))
//	return
//}

//func (m *httpHandler) getWxUserInfo(spId string, wxCode string) (errMsg string, openId string, sessionKey string, err error) {
//	var appid string
//	var secret string
//	tWxPluginProgram, err := m.wxpluginprogramsv.queryWxpluginProgram(spId)
//	if err != nil {
//		zap.L().Error(fmt.Sprintf("getWxUserInfo error %s", err.Error()))
//		return "", "", "", err
//	} else if tWxPluginProgram == nil {
//		return "系统未配置对应的小程序", "", "", nil
//	}
//	appid = tWxPluginProgram.Appid
//	secret = tWxPluginProgram.Appsecrete
//	resp, err := http.Get("https://api.weixin.qq.com/sns/jscode2session?appid=" + appid + "&secret=" + secret + "&js_code=" + wxCode + "&grant_type=authorization_code")
//	if err != nil {
//		zap.L().Error(fmt.Sprintf("get wx session_key error %s", err.Error()))
//		return "", "", "", err
//	}
//	body, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		zap.L().Error(fmt.Sprintf("get wx session_key error %s", err.Error()))
//		return "", "", "", err
//	}
//	var dat map[string]interface{}
//	if err := json.Unmarshal([]byte(body), &dat); err == nil {
//		openid := dat["openid"]
//		session_key := dat["session_key"]
//		errcode := dat["errcode"]
//		errMsg := dat["errMsg"]
//		zap.L().Debug(fmt.Sprintf("openid:%s", openid.(string)))
//		zap.L().Debug(fmt.Sprintf("session_key:%s", session_key.(string)))
//		zap.L().Debug(fmt.Sprintf("errcode:%s", errcode.(string)))
//		zap.L().Debug(fmt.Sprintf("errMsg:%s", errMsg.(string)))
//		if errcode.(int) == 0 {
//			return "", openid.(string), session_key.(string), nil
//		}
//		return "[" + wxCode + "] code2Session失败[" + strconv.Itoa(errcode.(int)) + "]:" + errMsg.(string), "", "", nil
//	}
//	return "", "", "", err
//}

//func (m *httpHandler) queryUser(body []byte, w http.ResponseWriter) {
//	var req QueryUserReq
//	err := json.Unmarshal(body, &req)
//	if err != nil {
//		zap.L().Error(fmt.Sprintf("json transfer error %s", err.Error()))
//		m.ivalidResp(w)
//		return
//	}
//	errMsg, openid, _, err := m.getWxUserInfo(req.Data.SpId, req.Data.WxCode)
//	if err != nil {
//		zap.L().Error(fmt.Sprintf("queryUser error %s", err.Error()))
//		m.ivalidResp(w)
//		return
//	} else if errMsg != "" {
//		var respHead ResponseHead
//		respHead = ResponseHead{RequestId: req.RequestHead.RequestId, ErrorCode: 9999, ErrorMsg: errMsg, Cmd: req.RequestHead.Cmd}
//		jsonData, err := json.Marshal(respHead)
//		if err != nil {
//			zap.L().Error(fmt.Sprintf("json transfer error %s", err.Error()))
//			m.ivalidResp(w)
//			return
//		}
//		w.WriteHeader(http.StatusOK)
//		w.Write([]byte(jsonData))
//		return
//	}
//	tUser, err := m.usersv.queryUserByOpenId(openid)
//	var resp QueryUserResp
//	if err == nil {
//		if tUser != nil {
//			var respHead ResponseHead
//			respHead = ResponseHead{RequestId: req.RequestHead.RequestId, ErrorCode: 0, Cmd: 131}
//			var data QueryUserRespData
//			data = QueryUserRespData{OpenId: tUser.Open_id.String, UserId: tUser.User_uuid, UserType: tUser.User_type, UserName: tUser.User_name.String, HeadIco: tUser.Head_portrait.String}
//			resp = QueryUserResp{ResponseHead: respHead, Data: data}
//		} else {
//			resp = QueryUserResp{ResponseHead{RequestId: req.RequestHead.RequestId, ErrorCode: 1, Cmd: 131, ErrorMsg: "未查到对应的用户"}, QueryUserRespData{}}
//		}
//		jsonData, err := json.Marshal(resp)
//		if err != nil {
//			zap.L().Error(fmt.Sprintf("json transfer error %s", err.Error()))
//			m.ivalidResp(w)
//			return
//		}
//		w.WriteHeader(http.StatusOK)
//		w.Write([]byte(jsonData))
//		return
//	}
//	zap.L().Error(fmt.Sprintf("queryUser error %s", err.Error()))
//	m.ivalidResp(w)
//}
