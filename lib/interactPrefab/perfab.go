package interactprefab

import (
	"goMiraiQQBot/lib/constdata"
	messagetargets "goMiraiQQBot/lib/messageHandle/messageTargets"
	"goMiraiQQBot/lib/messageHandle/sourceHandle"
	"goMiraiQQBot/lib/messageHandle/structs"
)

// 交互器预制件

type InteractPerfab struct {
	initFn func()

	useage string

	activateSigns  []string
	activateSource []constdata.MessageType
}




func (i *InteractPerfab) Init() {
	i.initFn()
}

//GetUseage
func (i *InteractPerfab) GetUseage() string {
	return i.useage
}

// GetCommandName 获取指令名称，用来创建映射关系
func (i *InteractPerfab) GetCommandName() []string {
	return i.activateSigns
}

//RespondSource
func (i *InteractPerfab) RespondSource() []constdata.MessageType {
	return i.activateSource
}

func(i*InteractPerfab)TodoMsgRespond(
	source sourceHandle.MessageSource,
	msgRsp chan messagetargets.MessageTarget){
	msgRsp<-messagetargets.SourceTarget(source,structs.NewTextChain("[预制件默认信息]本功能未完成，敬请期待"))
}