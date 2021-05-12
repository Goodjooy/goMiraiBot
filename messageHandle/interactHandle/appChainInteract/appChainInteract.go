package appchaininteract

import (
	"fmt"
	"goMiraiQQBot/constdata"
	datautil "goMiraiQQBot/dataUtil"
	"goMiraiQQBot/messageHandle/interact"
	messagetargets "goMiraiQQBot/messageHandle/messageTargets"
	"goMiraiQQBot/messageHandle/sourceHandle"
	"goMiraiQQBot/messageHandle/structs"
	"log"
)

type AppChainInteract struct {
}

func NewAppChainInteract() interact.ChainTypeInteract {
	return &AppChainInteract{}
}

//GetUseage 获取使用
func (xml *AppChainInteract) GetUseage() string {
	return `分享XML类型的分享，将解析xml并返回`
}

//GetActivateTypes 获取激活的chain类型
func (xml *AppChainInteract) GetActivateType() []constdata.MessageDataType {
	return []constdata.MessageDataType{
		constdata.App,
	}
}

//GetActivateSource 获取激活的信息类型
func (xml *AppChainInteract) GetActivateSource() []constdata.MessageType {
	return []constdata.MessageType{
		constdata.GroupMessage,
		constdata.FriendMessage,
		constdata.TempMessage,
	}
}

// EnterMessage 响应信息
func (xml *AppChainInteract) EnterMessage(
	extraCmd datautil.MutliToOneMap,
	data structs.Message,
	repChan chan messagetargets.MessageTarget) {
	var source, _ = sourceHandle.GetGoupSoucreMessage(data.Source)
	appInfo, err := jsonLoader(data.ChainInfoList[0].Data["content"].(string))
	if err != nil {
		if err != nil {
			log.Print("failure to Transform Json |", err)
			repChan <- messagetargets.NewSingleTextGroupTarget(source.GroupId, fmt.Sprintf("failure to Transform XML|%v", err))
		}
	}
	bodyMessage:=appInfo.Meta.Detial.Title
	if len(bodyMessage)>64{
		bodyMessage=bodyMessage[:64]+"……"
	}

	respond := messagetargets.NewChainsGroupTarget(source.GroupId,
		//structs.NewImageChain(xmlMessage.Source.Icon),
		structs.NewTextChain(appInfo.AppTitle),
		structs.NewImageChain(removeURLParams(appInfo.Meta.Detial.Preview)),
		structs.NewTextChain(bodyMessage),
		structs.NewTextChain("\n---------------\n"),
		structs.NewTextChain(fmt.Sprintf("网址：\n%s", removeURLParams(appInfo.Meta.Detial.QqDocURL))),
	)

	repChan <- respond
}
