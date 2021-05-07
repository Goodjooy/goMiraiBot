package sinplemessageinteract

import (
	"goMiraiQQBot/constdata"
	msh "goMiraiQQBot/messageHandle/messageSourceHandles"
	messagetargets "goMiraiQQBot/messageHandle/messageTargets"
	"goMiraiQQBot/messageHandle/structs"
	"goMiraiQQBot/request"
)

type AboutInteract struct {
}

func NewAboutInteract() SingleMessageInteract {
	return AboutInteract{}
}

func (AboutInteract) GetCommandName() string {
	return "about"
}
func (AboutInteract) RespondSource() []constdata.MessageType {
	return []constdata.MessageType{
		constdata.GroupMessage,
	}
}

func (i AboutInteract) EnterMessage(
	data structs.Message,
	repChan chan messagetargets.MessageTarget) {

		var msg msh.GroupMessage = data.Source.GetMetaInformation().(msh.GroupMessage)
	
		var da=messagetargets.NewGroupTarget(msg.GroupId,[]request.H{
			{
				"type":string(constdata.Plain),
				"text":"ForzenStringBot 是由凊弦凝绝制作的以Mirai为框架的QQ萝卜子",
			}})

		repChan<-da
	}
