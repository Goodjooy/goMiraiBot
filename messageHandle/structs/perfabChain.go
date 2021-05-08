package structs

import (
	"goMiraiQQBot/constdata"
)


func NewTextChain(msg string) MessageChainInfo {
	text := NewChain(constdata.Plain)
	text.SetData("text", msg)

	return text
}

func NewImageChain(imgURL string)MessageChainInfo{
	img:=NewChain(constdata.Image)
	img.SetData("url",imgURL)

	return img
}

func NewAtChain(qq uint64, display string) MessageChainInfo {
	at := NewChain(constdata.At)
	at.SetData("target", qq)

	return at
}

func NewQuoteChain(msgId,senderQQ,groupId uint64,origin []interface{})MessageChainInfo{
	quote:=NewChain(constdata.Quote)

	quote.SetData("id",msgId)
	quote.SetData("senderId",senderQQ)
	quote.SetData("targetId",groupId)
	quote.SetData("groupId",groupId)

	return quote
}