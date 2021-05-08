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

/*
ContextMessageInteract 有上下文关系的信息交互部分,
能够提供连续的信息交互.优先级高于普通信息交互
*/
type ContextMessageInteract interface {
	//GetUseage 获取命令使用方法
	GetUseage()string
	//GetInitCommand
	GetInitCommand()[]string

	//InitMessage 上下文交互创建时使用,初始化数据，响应消息
	InitMessage(extraCmd datautil.MutliToOneMap,msg structs.Message,redChan chan messagetargets.MessageTarget)ContextInteract
	//NextMessage 向上下文提交信息
	NextMessage(msg structs.Message,redChan chan messagetargets.MessageTarget)
	
	//IsDone 判断该上下文是否已经完成
	IsDone()bool
}