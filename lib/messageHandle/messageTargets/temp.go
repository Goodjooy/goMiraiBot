package messagetargets

import (
	"fmt"
	"goMiraiQQBot/lib/constdata"
	"goMiraiQQBot/lib/messageHandle/structs"
	"goMiraiQQBot/lib/request"
	"goMiraiQQBot/lib/request/structs/message"
)

type TempTarget struct {
	data message.TempMessageRequest
}

func (t TempTarget) GetTargetPort() constdata.MessageSendPort {
	return constdata.TempSend
}

func (t TempTarget) GetSendContain(sessionKey string) interface{} {
	return t.data
}
func (t TempTarget) GetSendMessage() string {
	return fmt.Sprintf("Send Temp Message<GID: %v><QQ: %v>", t.data.Group, t.data.QQ)
}

func NewTempTarget(groupId, userId uint64, chains []request.H) MessageTarget {
	data := message.TempMessageRequest{
		Group: groupId,
		QQ:    userId,
		Clain: (chains),
	}
	return TempTarget{
		data: data,
	}
}
func NewChainsTempTarget(groupId, userId uint64, chains ...structs.MessageChainInfo) MessageTarget {
	var datas []request.H

	for _, v := range chains {
		datas = append(datas, v.ToMap())
	}
	return NewTempTarget(groupId, userId, datas)
}

func NewSingleTextTempTarget(groupId, userId uint64, text string) MessageTarget {
	return NewTempTarget(groupId, userId, []request.H{
		{
			"type": constdata.Plain.String(),
			"text": text,
		},
	})
}
