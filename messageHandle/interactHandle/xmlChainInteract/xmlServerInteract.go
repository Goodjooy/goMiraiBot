package xmlchaininteract

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

type XmlChainInteract struct {
}

func NewXmlChainInteract() interact.ChainTypeInteract {
	return &XmlChainInteract{}
}
func (*XmlChainInteract) Init() {

}

//GetUseage 获取使用
func (xml *XmlChainInteract) GetUseage() string {
	return `分享XML类型的分享，将解析xml并返回`
}

//GetActivateTypes 获取激活的chain类型
func (xml *XmlChainInteract) GetActivateType() []constdata.MessageDataType {
	return []constdata.MessageDataType{
		constdata.Xml,
	}
}

//GetActivateSource 获取激活的信息类型
func (xml *XmlChainInteract) GetActivateSource() []constdata.MessageType {
	return []constdata.MessageType{
		constdata.GroupMessage,
		constdata.FriendMessage,
		constdata.TempMessage,
	}
}

// EnterMessage 响应信息
func (xml *XmlChainInteract) EnterMessage(
	extraCmd datautil.MutliToOneMap,
	data structs.Message,
	repChan chan messagetargets.MessageTarget) {
	var source, _ = sourceHandle.GetGoupSoucreMessage(data.Source)
	xmlMessage, err := loadXML(data.ChainInfoList[0].Data["xml"].(string))
	if err != nil {
		log.Fatal("failure to Transform XML|", err)
		repChan <- messagetargets.NewSingleTextGroupTarget(source.GroupId, fmt.Sprintf("failure to Transform XML|%v", err))
	}

	respond := messagetargets.NewChainsGroupTarget(source.GroupId,
		//structs.NewImageChain(xmlMessage.Source.Icon),
		structs.NewTextChain(xmlMessage.Item.Title),
		structs.NewImageChain(xmlMessage.Item.Picture.Cover),
		structs.NewTextChain(xmlMessage.Item.Summary),
		structs.NewTextChain("\n---------------\n"),
		structs.NewTextChain(fmt.Sprintf("网址：\n%s", removeURLParams(xmlMessage.URL))),
	)

	repChan <- respond

}
