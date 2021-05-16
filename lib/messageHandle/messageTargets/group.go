package messagetargets

import (
	"fmt"
	"goMiraiQQBot/lib/constdata"
	"goMiraiQQBot/lib/messageHandle/structs"
	"goMiraiQQBot/lib/request"
	"goMiraiQQBot/lib/request/structs/message"
)

type GroupTarget struct {
	data message.GroupMessageRequest
}

func NewGroupTarget(targetId uint64, chains []request.H) MessageTarget {
	data := message.GroupMessageRequest{
		Target: targetId,
		Clain:  (chains),
	}
	return GroupTarget{
		data: data,
	}
}
func NewChainsGroupTarget(groupId uint64, chains ...structs.MessageChainInfo) MessageTarget {
	var datas []request.H

	for _, v := range chains {
		datas = append(datas, v.ToMap())
	}
	return NewGroupTarget(groupId, datas)
}

func NewSingleTextGroupTarget(grpoupId uint64, text string) MessageTarget {
	return NewGroupTarget(grpoupId, []request.H{
		{
			"type": constdata.Plain.String(),
			"text": text,
		},
	})
}

func (target GroupTarget) GetTargetPort() constdata.MessageSendPort {
	return constdata.GroupSend
}
func (target GroupTarget) GetSendContain(sessionKey string) interface{} {
	target.data.Session = sessionKey
	return target.data
}

func (target GroupTarget)GetSendMessage()string{
	return fmt.Sprintf("Send Group Message<GID:%v>",target.data.Target)
}
