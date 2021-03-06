package interact

import (
	datautil "goMiraiQQBot/lib/dataUtil"
	messagetargets "goMiraiQQBot/lib/messageHandle/messageTargets"
	"goMiraiQQBot/lib/messageHandle/structs"
)

type Interact interface {
	Init()
	// EnterMessage 响应信息
	EnterMessage(extraCmd datautil.MutliToOneMap, data structs.Message, repChan chan messagetargets.MessageTarget)
}

type ContextInteract interface {
	Init()
	/*InitMessage 	
		上下文交互创建时使用,初始化数据，响应消息
		要求根据传递的数据将不完全初始化的响应器初始化完成
		并做出响应（可选）
	*/
	InitMessage(extraCmd datautil.MutliToOneMap, msg structs.Message, redChan chan messagetargets.MessageTarget) ContextInteract
	//NextMessage 向上下文提交信息
	NextMessage(msg structs.Message, redChan chan messagetargets.MessageTarget)

	//IsDone 判断该上下文是否已经完成
	IsDone() bool
}
