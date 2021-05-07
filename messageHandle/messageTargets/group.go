package messagetargets

import (
	"goMiraiQQBot/constdata"
	"goMiraiQQBot/request"
	"goMiraiQQBot/request/structs/message"
)


type GroupTarget struct {
	data message.GroupMessageRequest
}

func NewGroupTarget(targetId uint64,chains[]request.H) MessageTarget {
	data:=message.GroupMessageRequest{
		Target: targetId,
		Clain:  (chains),
	}
	return GroupTarget{
		data: data,
	}
}

func NewSingleTextGroupTarget(grpoupId uint64,text string)MessageTarget{
	return NewGroupTarget(grpoupId,[]request.H{
		{
			"type":constdata.Plain.String(),
			"text":text,
		},
	})
}

func (target GroupTarget) GetTargetPort() constdata.MessageSendPort {
	return constdata.GroupSend
}
func (target GroupTarget) GetSendContain(sessionKey string) interface{} {
	target.data.Session=sessionKey
	return target.data
}