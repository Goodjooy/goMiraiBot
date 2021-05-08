package messagetargets

import "goMiraiQQBot/constdata"

type MessageTarget interface {
	//GetTargetPort 获取信息发送的端口
	GetTargetPort()constdata.MessageSendPort
	//GetSendContain 获取荷载
	GetSendContain(sessionKey string)interface{}
	//GetSendMessage
	GetSendMessage()string
}

