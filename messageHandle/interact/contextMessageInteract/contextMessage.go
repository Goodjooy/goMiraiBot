package contextmessageinteract

import "goMiraiQQBot/messageHandle/structs"

/*
ContextMessageInteract 有上下文关系的信息交互部分,
能够提供连续的信息交互.优先级高于普通信息交互
*/
type ContextMessageInteract interface {
	//GetInitCommand
	GetInitCommand()string

	//InitMessage 上下文交互创建时使用,初始化数据，响应消息
	InitMessage(msg structs.Message,redChan chan structs.MessageRespond)ContextMessageInteract
	//NextMessage 向上下文提交信息
	NextMessage(msg structs.Message,redChan chan structs.MessageRespond)
	
	//IsDone 判断该上下文是否已经完成
	IsDone()bool
}