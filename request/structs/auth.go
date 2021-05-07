package structs

import "goMiraiQQBot/request"

type AuthRespond struct {
	Code    request.RespondStatus `json:"code"`
	Session string                `json:"session"`
}

type VerifyQQRequest struct {
	SessionKey string `json:"sessionKey"`
	QQ         uint64 `json:"qq"`
}

type VerifyRespond struct {
	Code    request.RespondStatus `json:"code"`
	Message string                `json:"msg"`
}
