package messagetargets

import (
	"fmt"
	"goMiraiQQBot/lib/constdata"
	"goMiraiQQBot/lib/messageHandle/structs"
	"goMiraiQQBot/lib/request"
	"goMiraiQQBot/lib/request/structs/message"
)


type FriendMessageTarget struct {
	data message.FriendMessageRequest
}

func NewFriendTarget(targetId uint64, chains []request.H) MessageTarget {
	data := message.FriendMessageRequest{
		QQ: targetId,
		Clain:  (chains),
	}
	return FriendMessageTarget{
		data: data,
	}
}
func NewChainsFriendTarget(qq uint64, chains ...structs.MessageChainInfo) MessageTarget {
	var datas []request.H

	for _, v := range chains {
		datas = append(datas, v.ToMap())
	}
	return NewFriendTarget(qq, datas)
}

func NewSingleTextFriendTarget(qq uint64, text string) MessageTarget {
	return NewChainsFriendTarget(qq,structs.NewTextChain(text))
}
//GetTargetPort 获取信息发送的端口
func (FriendMessageTarget)GetTargetPort()constdata.MessageSendPort{
	return constdata.FirendSend
}
//GetSendContain 获取荷载
func (target FriendMessageTarget)GetSendContain(sessionKey string)interface{}{
	target.data.Session=sessionKey
	return target.data
}
//GetSendMessage
func (target FriendMessageTarget)GetSendMessage()string{
	return fmt.Sprintf("Send Friend Message<FID:%v>",target.data.QQ)
}