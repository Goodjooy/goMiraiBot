package message

import (
	"goMiraiQQBot/constdata"
	"goMiraiQQBot/request"
)

type GroupMessageRequest struct {
	Session string `json:"sessionKey"`

	Target uint64 `json:"target"`
	Quote  uint64 `json:"-"`

	Clain []request.H `json:"messageChain"`
}
type TempMessageRequest struct {
	Session string `json:"sessionKey"`

	QQ    uint64 `json:"target"`
	Group uint64 `json:"-"`

	Clain []request.H `json:"messageChain"`
}

type MessageSendRespond struct {
	Code constdata.RespondStatus `json:"code"`

	Message   string `json:"msg"`
	MessageId uint64 `json:"messageId"`
}
