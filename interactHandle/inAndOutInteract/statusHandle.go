package inandoutinteract

import (
	"fmt"
	"goMiraiQQBot/messageHandle/structs"
	"regexp"
	"strconv"
	"time"

	"gorm.io/gorm"
)

var paymentPattern = regexp.MustCompile(`^(\S+)\s+?([＋+-—]?\d+)(?:\.(\d{0,2}))?\s*([^\s\.]?)$`)

func (i *PaymentRecordInteract) onAdd(
	msg structs.Message) []structs.MessageChainInfo {
	if i.status != ADD {
		return []structs.MessageChainInfo{}
	}

	if checkSingleAtCmd(msg.ChainInfoList) {
		i.status = Nil
		return []structs.MessageChainInfo{
			structs.NewTextChain("退出记账模式！"),
		}
	}

	mssage := cmdLoad(msg.ChainInfoList)

	data := paymentPattern.FindStringSubmatch(mssage)

	if data == nil {
		return []structs.MessageChainInfo{
			structs.NewTextChain("指令格式错误：请按照[账单信息] [收支金额] 进行记录，小数点后保留2位"),
		}
	}
	floatNum, err := strconv.Atoi(data[3])
	if err != nil {
		floatNum = 0
	}
	intNum, err := strconv.Atoi(data[2])
	if err != nil {
		intNum = 0
	}
	var goodType string
	if data[4] == "" {
		goodType = "￥"
	} else {
		goodType = data[4]
	}

	pay := PaymentRecord{
		Date:         time.Now(),
		PaymentFloat: uint16(floatNum),
		PaymentInt:   int32(intNum),
		Message:      data[1],
		GoodType:     goodType,
		UserId:       i.user.ID,
	}
	i.db.Save(&pay)

	return []structs.MessageChainInfo{
		structs.NewTextChain(fmt.Sprintf("保存账单成功！ \n %v.%v %v | %v",
			intNum, floatNum, goodType, data[1])),
	}
}

func (i *PaymentRecordInteract) onShowRecord(
	msg structs.Message) []structs.MessageChainInfo {

	var pays = loadPayments(i.user, i.db)

	var chains []structs.MessageChainInfo

	
	chains = append(chains, structs.NewTextChain("\n------\n"))

	var sum float64 = 0
	for _, v := range pays {
		chains = append(chains,
			structs.NewTextChain(fmt.Sprintf("%v.%v %v\n  %v\n  %v\n-------\n",
				v.PaymentInt,
				v.PaymentFloat,
				v.GoodType, v.Message,
				v.Date.Format("2006.01.02 15:04:05"),
			)))
		var floatPart float64
		if v.PaymentFloat == 0 {
			floatPart = 0
		} else {
			floatPart = float64(v.PaymentFloat) / 100
		}
		sum += float64(v.PaymentInt) + floatPart
	}
	chains = append(chains, structs.NewTextChain(fmt.Sprintf("总计： %.2f", sum)))

	return chains
}

func (i *PaymentRecordInteract) onNil(
	msg structs.Message) []structs.MessageChainInfo {
	message := cmdLoad(msg.ChainInfoList)
	if message == exit {
		i.status = EXIT
		return []structs.MessageChainInfo{
			structs.NewTextChain("关闭记账"),
		}
	} else if message == loadPayment {
		chains := i.onShowRecord(msg)
		return chains
	} else if message == addPayment {
		i.status = ADD
		return []structs.MessageChainInfo{
			structs.NewTextChain("开始记账，\n格式： [账单信息] [收支金额],\n@我或者‘结束记账’终止"),
		}
	} else {
		return []structs.MessageChainInfo{
			structs.NewTextChain("未知指令！"),
		}
	}
}

func loadPayments(user User, db *gorm.DB) []PaymentRecord {
	var pay PaymentRecord = PaymentRecord{
		UserId: user.ID,
		User:   user,
	}
	var pays []PaymentRecord
	db.Where(&pay).Find(&pays)

	return pays
}
