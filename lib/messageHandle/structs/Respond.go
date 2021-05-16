package structs

import "goMiraiQQBot/lib/constdata"

type MessageRespond struct {
	Port          constdata.MessageSendPort
	MessageChains []MessageChainInfo
}