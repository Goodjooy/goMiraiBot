package sourceHandle

import (
	"errors"
	"fmt"
	"goMiraiQQBot/constdata"
	"goMiraiQQBot/request/structs/message"
	"log"
)

type MessageSource interface {
	// GetSource 获取消息的来源
	GetSource() constdata.MessageType
	//GetMetaInformation 获取信息来源的内部信息,即sender信息
	GetMetaInformation() interface{}

	//GetSenderUserID get send user qq id
	GetSenderID() uint64
	//GetGroupID get message form group ,firend return 0
	GetGroupID() uint64
}

var sourceMap map[constdata.MessageType]func(map[string]interface{}) MessageSource = make(map[constdata.MessageType]func(map[string]interface{}) MessageSource)

func InitMessageSourceHandle() {
	sourceMap[constdata.GroupMessage] = FromMessageRecive
	sourceMap[constdata.FriendMessage] = FromMessageReciveToFriend
}

func FromMessageMap(data message.MessageMapRespond) (MessageSource, error) {
	if data.Code != constdata.Normal {
		return nil, errors.New("Failure Operate ErrorCode : " + fmt.Sprint(data.Code) + "Info: " + data.ErrorMessage)
	}

	defer handingMessagePanicReover()

	metaData := data.Data
	var messageType = constdata.MessageType(metaData["type"].(string))
	handle,ok:= sourceMap[messageType]
	if !ok{
		return nil,fmt.Errorf("no Message Handle For %v Found",messageType)
	}

	var messageSource = handle(metaData["sender"].(map[string]interface{}))

	return messageSource, nil
}

func handingMessagePanicReover() {
	err := recover()
	if err == nil {
		log.Print("Handing message well")
	} else {
		log.Print("HandingMessage Error: ", err)
	}
}
