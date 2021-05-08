package structs

import (
	"goMiraiQQBot/constdata"
)

type MessageChainInfo struct {
	MessageType constdata.MessageDataType
	Data        map[string]interface{}
}

func NewChain(chainType constdata.MessageDataType) MessageChainInfo {
	return MessageChainInfo{MessageType: chainType,Data: make(map[string]interface{})}
}


func (mci *MessageChainInfo) ToMap() map[string]interface{} {
	mci.Data["type"] = mci.MessageType.String()

	return mci.Data
}

func (mci *MessageChainInfo) SetData(key string, value interface{}) {
	mci.Data[key] = value
}
