package messagetargets

import (
	"goMiraiQQBot/constdata"
	"goMiraiQQBot/messageHandle/sourceHandle"
	"goMiraiQQBot/messageHandle/structs"
	"log"
)

type MessageTarget interface {
	//GetTargetPort 获取信息发送的端口
	GetTargetPort() constdata.MessageSendPort
	//GetSendContain 获取荷载
	GetSendContain(sessionKey string) interface{}
	//GetSendMessage
	GetSendMessage() string
}

func SourceTarget(source sourceHandle.MessageSource, chains ...structs.MessageChainInfo) MessageTarget {
	if source.GetSource() == constdata.GroupMessage {
		return NewChainsGroupTarget(source.GetGroupID(), chains...)
	} else if source.GetSource() == constdata.FriendMessage {
		return NewChainsFriendTarget(source.GetSenderID(), chains...)
	} else if source.GetSource() == constdata.TempMessage {
		return NewChainsTempTarget(source.GetGroupID(), source.GetSenderID(), chains...)
	}
	log.Printf("Same Way Repeat Failure: unknow Message Type: %v", source.GetSource())
	return nil
}
