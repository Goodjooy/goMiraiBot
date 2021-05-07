package messagesourcehandles

import (
	"goMiraiQQBot/constdata"
)

type GroupMessage struct {
	UserId   uint64
	UserName string

	GroupId   uint64
	GroupName string

	SenderPermission constdata.PermissionLevel
	BotPermission    constdata.PermissionLevel
}

type GroupMessageSourceHandle struct {
	data GroupMessage
}

func FromMessageRecive(sender map[string]interface{}) MessageSource {
	var group = sender["group"].(map[string]interface{})

	data := GroupMessage{
		UserId:           uint64(sender["id"].(float64)),
		UserName:         sender["memberName"].(string),
		SenderPermission: constdata.PermissionLevel(sender["permission"].(string)),
		GroupId:          uint64(group["id"].(float64)),
		GroupName:        group["name"].(string),
		BotPermission:    constdata.PermissionLevel(group["permission"].(string)),
	}
	return GroupMessageSourceHandle{
		data: data,
	}

}

func (handle GroupMessageSourceHandle) GetSource() constdata.MessageType {
	return constdata.GroupMessage
}

func (handle GroupMessageSourceHandle) GetMetaInformation() interface{} {
	return handle.data
}
