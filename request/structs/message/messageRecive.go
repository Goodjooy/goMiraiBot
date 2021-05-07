package message

import "goMiraiQQBot/request"

type MessageRespond struct {
	Code         request.RespondStatus `json:"code"`
	ErrorMessage string                `json:"errorMessage"`

	Data MessageData `json:"data"`
}

type MessageData struct {
	Type request.MessageType `json:"type"`

	MessageChain []ChainInfo `json:"messageChain"`
	Sender       Sender       `json:"sender"`
}

type ChainInfo map[string]interface{}

type Sender struct {
	Id         uint64 `json:"id"`
	MemberName string `json:"memberName"`


	Permission request.PermissionLevel `json:"permission"`

	GroupIn Group `json:"group"`
}

type Group struct {
	Id   uint64 `json:"id"`
	Name string`json:"name"`

	Permission request.PermissionLevel `json:"permission"`
}
