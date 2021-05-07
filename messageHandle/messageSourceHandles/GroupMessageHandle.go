package messagesourcehandles

import (
	"goMiraiQQBot/constdata"
)

type GroupMessage struct {
	userId   uint64
	userName string

	groupId   uint64
	groupName string

	senderPermission constdata.PermissionLevel
	botPermission    constdata.PermissionLevel
}

type GroupMessageSourceHandle struct {
	data GroupMessage
}

func FromMessageRecive(sender map[string]interface{}) MessageSource {
	var group = sender["group"].(map[string]interface{})

	data := GroupMessage{
		userId:           sender["id"].(uint64),
		userName:         sender["memberName"].(string),
		senderPermission: constdata.PermissionLevel(sender["permission"].(string)),
		groupId:          uint64(group["id"].(float64)),
		groupName:        group["name"].(string),
		botPermission:    constdata.PermissionLevel(group["permission"].(string)),
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
