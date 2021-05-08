package hentaiimageinteract

import (
	datautil "goMiraiQQBot/dataUtil"
	"goMiraiQQBot/messageHandle/interact"
	messagetargets "goMiraiQQBot/messageHandle/messageTargets"
	"goMiraiQQBot/messageHandle/sourceHandle"
	"goMiraiQQBot/messageHandle/structs"
	"strconv"
)

type ContextInteract interact.ContextMessageInteract

type HentaiImageSearchInteract struct {
	imageUrl   string
	sendNumber int
	done       bool

	groupId uint64
	userId  uint64
}

func NewHentaiImageSearchInteract() interact.ContextMessageInteract {
	return &HentaiImageSearchInteract{done: false, imageUrl: ""}
}

func (HentaiImageSearchInteract) GetInitCommand() []string {
	return []string{"s-Img", "搜图"}
}

func (h *HentaiImageSearchInteract) InitMessage(
	extraCmd datautil.MutliToOneMap,
	data structs.Message,
	redChan chan messagetargets.MessageTarget) interact.ContextInteract {
	var msg sourceHandle.GroupMessage = data.Source.GetMetaInformation().(sourceHandle.GroupMessage)
	var res = messagetargets.NewChainsGroupTarget(msg.GroupId,
		structs.NewTextChain("请发送一个图片以搜索 \n 发送 “取消” 以取消等待图片操作"))
	redChan <- res

	countS, _ := extraCmd.GetWithDefault("3", "count", "长度", "数量")

	count, err := strconv.Atoi(countS)
	if err == nil {
		h.sendNumber = int(count)
	} else {
		h.sendNumber = 3
	}

	h.groupId = msg.GroupId
	h.userId = msg.UserId

	return h
}

func (h *HentaiImageSearchInteract) NextMessage(
	data structs.Message,
	redChan chan messagetargets.MessageTarget) {

	imgURL, ok := foundTargeImage(data.ChainInfoList)
	if ok {
		h.done = true
		go searchHandle(h.sendNumber, h.groupId, h.userId, imgURL, redChan)
	} else {
		if findCancelSign(data.ChainInfoList) {
			//取消信息接受
			h.done = true
			var m = messagetargets.NewSingleTextGroupTarget(h.groupId, "搜图任务取消")
			redChan <- m
			return
		} else {
			var m = messagetargets.NewSingleTextGroupTarget(h.groupId, "错误指令，请发送一张图片")
			redChan <- m
			return
		}
	}
}

func (h *HentaiImageSearchInteract) IsDone() bool {
	return h.done
}

func (h *HentaiImageSearchInteract) GetUseage() string {
	return `#s-Img|#搜图 ： 通过给定图片进行相似图片查询
	额外指令：
		count|长度|数量=[大于0的整数] -> 发送的搜索结果数量，如果大于结果或者小于0就发送全部结果，默认：3`
}