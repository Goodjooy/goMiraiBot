package message

import (
	"goMiraiQQBot/constdata"
	"goMiraiQQBot/request"
)


type GroupMessageRequest struct {
	Session string `json:"sessionKey"`

	Target uint64 `json:"target"`

	Clain []request.H `json:"messageChain"`
}

type MessageSendRespond struct {
	Code constdata.RespondStatus `json:"code"`

	Message   string `json:"msg"`
	MessageId uint64 `json:"messageId"`
}
