package structs

import (
	"goMiraiQQBot/constdata"
	messagesourcehandles "goMiraiQQBot/messageHandle/messageSourceHandles"
	"goMiraiQQBot/request/structs/message"
)

// 接受到的信息
type Message struct {
	//Source 信息的来源
	Source messagesourcehandles.MessageSource
	// HeadInfo 信息的头部信息,包括信息id和信息发送时间
	HeadInfo MessageHeadInfo
	//信息链的信息
	ChainInfoList []MessageChainInfo
}

type MessageHeadInfo struct {
	Id   uint64
	Time uint64
}

type MessageChainInfo struct {
	MessageType constdata.MessageDataType
	Data        map[string]interface{}
}

func FromMessageRespondData(data message.MessageMapRespond) (Message, error) {
	source, err := messagesourcehandles.FromMessageMap(data)
	if err != nil {
		return Message{}, err
	}
	dataMessage := data.Data
	messageChains := dataMessage["messageChain"].([]interface{})

	//headMessage
	sourceChain := messageChains[0].(map[string]interface{})
	head := MessageHeadInfo{
		Id:   uint64(sourceChain["id"].(float64)),
		Time: uint64(sourceChain["time"].(float64)),
	}

	var msg = Message{
		Source:   source,
		HeadInfo: head,
	}
	//messageChain
	for _, v := range messageChains[1:] {
		var metaInfo MessageChainInfo
		messageChainMetaData := v.(map[string]interface{})

		metaInfo.MessageType = constdata.MessageDataType(messageChainMetaData["type"].(string))

		metaInfo.Data=make(map[string]interface{})
		for k, d := range messageChainMetaData {
			if k == "type" {
				continue
			}
			metaInfo.Data[k] = d
		}
		msg.ChainInfoList = append(msg.ChainInfoList, metaInfo)
	}
	return msg, nil
}
