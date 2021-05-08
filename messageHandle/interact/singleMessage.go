package interact

import (
	"goMiraiQQBot/constdata"
	datautil "goMiraiQQBot/dataUtil"
	messagetargets "goMiraiQQBot/messageHandle/messageTargets"
	"goMiraiQQBot/messageHandle/structs"
)

type SingleMessageInteract interface {
	// GetCommandName 获取指令名称，用来创建映射关系
	GetCommandName() []string
	//GetUseage
	GetUseage() string

	//RespondSource
	RespondSource() []constdata.MessageType
	// EnterMessage 响应信息
	EnterMessage(extraCmd datautil.MutliToOneMap, data structs.Message, repChan chan messagetargets.MessageTarget)
}
