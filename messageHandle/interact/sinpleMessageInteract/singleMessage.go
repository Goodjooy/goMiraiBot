package sinplemessageinteract

import "goMiraiQQBot/messageHandle/structs"

type SingleMessageInteract interface {
	// GetCommandName 获取指令名称，用来创建映射关系
	GetCommandName()string
	// EnterMessage 响应信息
	EnterMessage(data structs.Message,repChan chan structs.MessageRespond)
}