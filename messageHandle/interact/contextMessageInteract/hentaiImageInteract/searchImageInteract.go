package hentaiimageinteract

import (
	"goMiraiQQBot/constdata"
	contextmessageinteract "goMiraiQQBot/messageHandle/interact/contextMessageInteract"
	messagesourcehandles "goMiraiQQBot/messageHandle/messageSourceHandles"
	messagetargets "goMiraiQQBot/messageHandle/messageTargets"
	"goMiraiQQBot/messageHandle/structs"
	"goMiraiQQBot/request"
)

type ContextInteract contextmessageinteract.ContextMessageInteract

type HentaiImageSearchInteract struct {
	imageUrl string
	done     bool
}

func NewHentaiImageSearchInteract() contextmessageinteract.ContextMessageInteract {
	return &HentaiImageSearchInteract{done: false, imageUrl: ""}
}

func (HentaiImageSearchInteract) GetInitCommand() string {
	return "搜图"
}

func (h *HentaiImageSearchInteract) InitMessage(
	data structs.Message,
	redChan chan messagetargets.MessageTarget) contextmessageinteract.ContextMessageInteract {
	var msg messagesourcehandles.GroupMessage = data.Source.GetMetaInformation().(messagesourcehandles.GroupMessage)
	var res = messagetargets.NewGroupTarget(msg.GroupId, []request.H{
		{
			"type": string(constdata.Plain),
			"text": "请发送一个图片以搜索 | 发送 “取消” 以取消等待图片操作",
		}})

	redChan <- res

	return h
}

func (h *HentaiImageSearchInteract) NextMessage(
	data structs.Message,
	redChan chan messagetargets.MessageTarget) {
	var msg messagesourcehandles.GroupMessage = data.Source.GetMetaInformation().(messagesourcehandles.GroupMessage)

	imgURL, ok := foundTargeImage(data.ChainInfoList)
	if ok {
		h.done=true
		go searchHandle(msg.GroupId,imgURL,redChan)
	} else {
		if findCancelSign(data.ChainInfoList) {
			//取消信息接受
			h.done = true
			var m = messagetargets.NewSingleTextGroupTarget(msg.GroupId, "搜图任务取消")
			redChan <- m
			return
		}else {
			var m = messagetargets.NewSingleTextGroupTarget(msg.GroupId, "错误指令，请发送一张图片")
			redChan <- m
			return
		}
	}
}

func (h *HentaiImageSearchInteract) IsDone() bool {
	return h.done
}
