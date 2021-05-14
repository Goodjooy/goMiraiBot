package interactHandle

import (
	"goMiraiQQBot/constdata"
	datautil "goMiraiQQBot/dataUtil"
	"goMiraiQQBot/messageHandle/interact"
	messagetargets "goMiraiQQBot/messageHandle/messageTargets"
	"goMiraiQQBot/messageHandle/sourceHandle"
	"goMiraiQQBot/messageHandle/structs"
	"goMiraiQQBot/request"
)

type AboutInteract struct {
}

func NewAboutInteract() interact.FullSingleInteract {
	return AboutInteract{}
}
func (AboutInteract) Init() {

}

func (AboutInteract) GetCommandName() []string {
	return []string{"about", "关于"}
}
func (AboutInteract) RespondSource() []constdata.MessageType {
	return []constdata.MessageType{
		constdata.GroupMessage,
	}
}

func (i AboutInteract) EnterMessage(
	extraCmd datautil.MutliToOneMap,
	data structs.Message,
	repChan chan messagetargets.MessageTarget) {

	var msg, _ = sourceHandle.GetGoupSoucreMessage(data.Source)

	var da = messagetargets.NewGroupTarget(msg.GroupId, []request.H{
		{
			"type": string(constdata.Plain),
			"text": "ForzenStringBot 是由凊弦凝绝制作的以Mirai为框架的QQ萝卜子",
		}})

	repChan <- da
}

func (AboutInteract) GetUseage() string {
	return "#about|#关于 : 返回当前机器人的简介\n" +
		"额外指令：无"
}
