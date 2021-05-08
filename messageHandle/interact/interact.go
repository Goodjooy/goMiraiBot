package interact

import (
	datautil "goMiraiQQBot/dataUtil"
	messagetargets "goMiraiQQBot/messageHandle/messageTargets"
	"goMiraiQQBot/messageHandle/structs"
)

type Interact interface {
	// EnterMessage 响应信息
	EnterMessage(extraCmd datautil.MutliToOneMap, data structs.Message, repChan chan messagetargets.MessageTarget)
}

type ContextInteract interface {
	//InitMessage 上下文交互创建时使用,初始化数据，响应消息
	InitMessage(extraCmd datautil.MutliToOneMap, msg structs.Message, redChan chan messagetargets.MessageTarget) ContextInteract
	//NextMessage 向上下文提交信息
	NextMessage(msg structs.Message, redChan chan messagetargets.MessageTarget)

	//IsDone 判断该上下文是否已经完成
	IsDone() bool
}
