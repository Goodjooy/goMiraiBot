package inandoutinteract

import (
	"goMiraiQQBot/constdata"
	"goMiraiQQBot/messageHandle/structs"
)

type Status int

const (
	Nil Status = iota
	SHOW
	ADD
	EXIT
)
const (
	loadPayment = "查看账单"
	addPayment  = "开始记账"
	exit        = "退出"
)

const nilMsgText = `输入‘开始记账’开始记录账单，
输入‘查看账单’查看账单,
输入‘退出’结束`

func cmdLoad(chains []structs.MessageChainInfo) string {
	var cmds string = ""

	for _, v := range chains {
		if v.MessageType == constdata.Plain {
			cmds += v.Data["text"].(string)
		}
	}
	return cmds
}

func checkSingleAtCmd(chains []structs.MessageChainInfo) bool {
	if len(chains) >= 1 && chains[0].MessageType == constdata.At {
		chain := chains[0]
		qq := uint64(chain.Data["target"].(float64))
		if qq == botQQ {
			return true
		}
	}else if  len(chains) >= 1 &&chains[0].MessageType==constdata.Plain {
		if chains[0].Data["text"]=="结束记账"{
			return true
		}
	}
	return false
}
