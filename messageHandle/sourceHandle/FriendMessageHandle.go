package sourceHandle

import "goMiraiQQBot/constdata"

type FriendMessage struct {
	UserId   uint64
	UserName string

	Remark string
}

type FriendMessageSourceHandle struct {
	data FriendMessage
}

func FromMessageReciveToFriend(sender map[string]interface{}) MessageSource {
	data := FriendMessage{
		UserId:   uint64(sender["id"].(float64)),
		UserName: string(sender["nickname"].(string)),
		Remark:   sender["remark"].(string),
	}
	return FriendMessageSourceHandle{
		data: data,
	}
}

func (handle FriendMessageSourceHandle) GetSource() constdata.MessageType {
	return constdata.FriendMessage
}

func (handle FriendMessageSourceHandle) GetMetaInformation() interface{} {
	return handle.data
}
func (handle FriendMessageSourceHandle) GetSenderID() uint64 {
	return handle.data.UserId
}
func (handle FriendMessageSourceHandle) GetGroupID() uint64 {
	return 0
}
