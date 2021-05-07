package structs

import "goMiraiQQBot/constdata"

type AuthRespond struct {
	Code    constdata.RespondStatus `json:"code"`
	Session string                `json:"session"`
}

type VerifyQQRequest struct {
	SessionKey string `json:"sessionKey"`
	QQ         uint64 `json:"qq"`
}

type VerifyRespond struct {
	Code    constdata.RespondStatus `json:"code"`
	Message string                `json:"msg"`
}
