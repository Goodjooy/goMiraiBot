package messagesourcehandles

import (
	"errors"
	"goMiraiQQBot/constdata"
	"goMiraiQQBot/request/structs/message"
)

type MessageSource interface {
	// GetSource 获取消息的来源
	GetSource() constdata.MessageType
	//GetMetaInformation 获取信息来源的内部信息,即sender信息
	GetMetaInformation() interface{}
}

var sourceMap map[constdata.MessageType]func(map[string]interface{}) MessageSource

func InitMessageSourceHandle() {
	sourceMap[constdata.GroupMessage] = FromMessageRecive
}

func FromMessageMap(data message.MessageMapRespond) (MessageSource, error) {
	if data.Code != constdata.Normal {
		return nil, errors.New("Failure Operate ErrorCode : " + string(data.Code) + "Info: " + data.ErrorMessage)
	}
	metaData := data.Data.(map[string]interface{})
	var messageType = constdata.MessageType(metaData["type"].(string))
	var handle = sourceMap[messageType]

	var messageSource = handle(metaData["sender"].(map[string]interface{}))

	return messageSource, nil
}
