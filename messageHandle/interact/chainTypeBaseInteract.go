package interact

import (
	"goMiraiQQBot/constdata"
	datautil "goMiraiQQBot/dataUtil"
	messagetargets "goMiraiQQBot/messageHandle/messageTargets"
	"goMiraiQQBot/messageHandle/structs"
)


type ChainTypeInteract interface {
	//GetUseage 获取使用
	GetUseage() string
	//GetActivateTypes 获取激活的chain类型
	GetActivateType()[]constdata.MessageDataType
	//GetActivateSource 获取激活的信息类型
	GetActivateSource()[]constdata.MessageType

	// EnterMessage 响应信息
	EnterMessage(extraCmd datautil.MutliToOneMap, data structs.Message, repChan chan messagetargets.MessageTarget)
}

type ChainTypeContextInteract interface {
	//GetUseage 获取命令使用方法
	GetUseage()string
	//GetActivateTypes 获取激活的chain类型
	GetActivateType()[]constdata.MessageDataType
	//GetActivateSource 获取激活的信息类型
	GetActivateSource()[]constdata.MessageType


	//InitMessage 上下文交互创建时使用,初始化数据，响应消息
	InitMessage(extraCmd datautil.MutliToOneMap,msg structs.Message,redChan chan messagetargets.MessageTarget)ChainTypeContextInteract
	//NextMessage 向上下文提交信息
	NextMessage(msg structs.Message,redChan chan messagetargets.MessageTarget)
	
	//IsDone 判断该上下文是否已经完成
	IsDone()bool
}