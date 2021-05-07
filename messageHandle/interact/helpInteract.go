package interact

import (
	"goMiraiQQBot/constdata"
	sinplemessageinteract "goMiraiQQBot/messageHandle/interact/sinpleMessageInteract"
	msh "goMiraiQQBot/messageHandle/messageSourceHandles"
	messagetargets "goMiraiQQBot/messageHandle/messageTargets"
	"goMiraiQQBot/messageHandle/structs"
	"goMiraiQQBot/request"
)

//帮助功能，用于显示全部命令
type HelpInteract struct {
}

func NewHelpInteract() sinplemessageinteract. SingleMessageInteract {
	return HelpInteract{}
}

func (interact HelpInteract) GetCommandName() string {
	return "help"
}
func (interact HelpInteract) RespondSource() []constdata.MessageType {
	return []constdata.MessageType{
		constdata.GroupMessage,
	}
}

func (i HelpInteract) EnterMessage(
	data structs.Message,
	repChan chan messagetargets.MessageTarget) {
	var msg msh.GroupMessage = data.Source.GetMetaInformation().(msh.GroupMessage)
	
	var d []request.H

	d=append(d, request.H{"type": string(constdata.Plain),
	"text": "收到信息：帮助\n"},)

	d=append(d, request.H{
		"type":string(constdata.Plain),
		"text":"以下为单次交互命令\n",
	})

	for _,v:=range GetSingleCommand(){
		d=append(d, request.H{
			"type":string(constdata.Plain),
			"text":v+"\n",
		})
	}

	d=append(d, request.H{
		"type":string(constdata.Plain),
		"text":"以下为上下文交互命令\n",
	})

	for _,v:=range GetContextCommand(){
		d=append(d, request.H{
			"type":string(constdata.Plain),
			"text":v+"\n",
		})
	}

	var da = messagetargets.NewGroupTarget(msg.GroupId, d)

	repChan <- da
}
