package interact
import (
	"goMiraiQQBot/lib/constdata"
	datautil "goMiraiQQBot/lib/dataUtil"
	messagetargets "goMiraiQQBot/lib/messageHandle/messageTargets"
	"goMiraiQQBot/lib/messageHandle/structs"
)


type FullSingleInteract interface {
	Init()
	//GetUseage
	GetUseage() string

	// GetCommandName 获取指令名称，用来创建映射关系
	GetCommandName() []string
	//RespondSource
	RespondSource() []constdata.MessageType

	// EnterMessage 响应信息
	EnterMessage(extraCmd datautil.MutliToOneMap, data structs.Message, repChan chan messagetargets.MessageTarget)
}

/*
FullContextInteract 有上下文关系的信息交互部分,
能够提供连续的信息交互.优先级高于普通信息交互
*/
type FullContextInteract interface {
	Init()
	//GetUseage 获取命令使用方法
	GetUseage()string
	//GetInitCommand
	GetCommandName()[]string
	//RespondSource
	RespondSource() []constdata.MessageType

	//InitMessage 上下文交互创建时使用,初始化数据，响应消息
	InitMessage(extraCmd datautil.MutliToOneMap,msg structs.Message,redChan chan messagetargets.MessageTarget)ContextInteract
	//NextMessage 向上下文提交信息
	NextMessage(msg structs.Message,redChan chan messagetargets.MessageTarget)
	
	//IsDone 判断该上下文是否已经完成
	IsDone()bool
}