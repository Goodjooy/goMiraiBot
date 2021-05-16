package message

import (
	"goMiraiQQBot/lib/constdata"
)

type MessageMapRespond struct {
	Code         constdata.RespondStatus `json:"code"`
	ErrorMessage string                `json:"errorMessage"`

	Data map[string]interface{} `json:"data"`
}
type MessageRespond struct {
	Code         constdata.RespondStatus `json:"code"`
	ErrorMessage string                `json:"errorMessage"`

	Data MessageData `json:"data"`
}

type MessageData struct {
	Type constdata.MessageType `json:"type"`

	MessageChain []ChainInfo `json:"messageChain"`
	Sender       Sender       `json:"sender"`
}

type ChainInfo map[string]interface{}

type Sender struct {
	Id         uint64 `json:"id"`
	MemberName string `json:"memberName"`


	Permission constdata.PermissionLevel `json:"permission"`

	GroupIn Group `json:"group"`
}

type Group struct {
	Id   uint64 `json:"id"`
	Name string`json:"name"`

	Permission constdata.PermissionLevel `json:"permission"`
}
