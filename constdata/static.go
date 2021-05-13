package constdata

type RespondStatus uint

const (
	Normal RespondStatus = iota
	BadAuthKey
	TargetBotNotExist
	SessionOutOfTimeOrNotExist
	SessionNotVerifyed
	MessageSendTargetNotExist
	TargetFileNotExist

	PemissionDenied           RespondStatus = 10
	BotCanNotSendMessage      RespondStatus = 20
	MessageSizeOverLimitation RespondStatus = 30
	BadRequest                RespondStatus = 400
)

type MessageType string

const (
	GroupMessage  MessageType = "GroupMessage"
	FriendMessage MessageType = "FriendMessage"
	TempMessage   MessageType = "TempMessage"
)

type PermissionLevel string

const (
	OWNER         PermissionLevel = "OWNER"
	ADMINISTRATOR PermissionLevel = "ADMINISTRATOR"
	MEMBER        PermissionLevel = "MEMBER"
)

type MessageDataType string

func (mdt MessageDataType) String() string { return string(mdt) }

const (
	Plain  MessageDataType = "Plain"
	Image  MessageDataType = "Image"
	Source MessageDataType = "Source"
	Quote  MessageDataType = "Quote"
	At     MessageDataType = "At"
	Xml    MessageDataType = "Xml"
	App    MessageDataType = "App"
)

type MessageSendPort string

const (
	GroupSend  MessageSendPort = "/sendGroupMessage"
	FirendSend MessageSendPort = "/sendFriendMessage"
	TempSend   MessageSendPort = "/sendTempMessage"
	ImageSend  MessageSendPort = "/sendImageMessage"

	RecallMessage MessageSendPort = "/recall"

	UploadImage MessageSendPort = "/uploadImage"
	UploadVoice MessageSendPort = "/uploadVoice"
	UploadFile  MessageSendPort = "/uploadFileAndSend"
)

type DatabaseName string

var (
	MYSQL  DatabaseName = "mysql"
	SQLITE DatabaseName = "sqlite"
)
