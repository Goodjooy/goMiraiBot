package structs

import "goMiraiQQBot/constdata"

type MessageRespond struct {
	Port          constdata.MessageSendPort
	MessageChains []MessageChainInfo
}