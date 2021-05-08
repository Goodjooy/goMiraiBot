package messagetargets

import (
	"goMiraiQQBot/constdata"
	"goMiraiQQBot/messageHandle/structs"
	"goMiraiQQBot/request"
	"goMiraiQQBot/request/structs/message"
)

type TempTarget struct {
	data message.TempMessageRequest
}

func (t TempTarget) GetTargetPort() constdata.MessageSendPort {
	return constdata.FirendSend
}

func (t TempTarget) GetSendContain(sessionKey string) interface{} {
	return t.data
}


func NewTempTarget(groupId,userId uint64, chains []request.H) MessageTarget {
	data := message.TempMessageRequest{
		Group: groupId,
		QQ: userId,
		Clain:  (chains),
	}
	return TempTarget{
		data: data,
	}
}
func NewChainsTempTarget(groupId,userId uint64, chains ...structs.MessageChainInfo) MessageTarget {
	var datas []request.H

	for _, v := range chains {
		datas = append(datas, v.ToMap())
	}
	return NewTempTarget(groupId,userId, datas)
}

func NewSingleTextTempTarget(groupId,userId uint64, text string) MessageTarget {
	return NewTempTarget(groupId,userId, []request.H{
		{
			"type": constdata.Plain.String(),
			"text": text,
		},
	})
}
