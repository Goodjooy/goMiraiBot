package inandoutinteract

import (
	"goMiraiQQBot/constdata"
	datautil "goMiraiQQBot/dataUtil"
	"goMiraiQQBot/database"
	"goMiraiQQBot/messageHandle/interact"
	messagetargets "goMiraiQQBot/messageHandle/messageTargets"
	"goMiraiQQBot/messageHandle/structs"

	"gorm.io/gorm"
)

var botQQ uint64

type PaymentRecordInteract struct {
	user User
	db   *gorm.DB

	status Status
}

func NewPaymentRecordInteract() interact.FullContextInteract {
	return &PaymentRecordInteract{}
}

func (i *PaymentRecordInteract) Init() {
	database.AsignDBModel(&PaymentRecord{}, &User{})
	botQQ = interact.GetBotQQ()
}

//GetUseage 获取命令使用方法
func (i *PaymentRecordInteract) GetUseage() string {
	return `#记账 - 上下文类型命令交互
		在之后可进行多种指令操作`
}

//GetInitCommand
func (i *PaymentRecordInteract) GetCommandName() []string {
	return []string{
		"记账", "payment-record",
	}
}

//RespondSource
func (i *PaymentRecordInteract) RespondSource() []constdata.MessageType {
	return []constdata.MessageType{
		constdata.GroupMessage,
		constdata.FriendMessage,
	}
}

//InitMessage 上下文交互创建时使用,初始化数据，响应消息
func (i *PaymentRecordInteract) InitMessage(
	extraCmd datautil.MutliToOneMap,
	msg structs.Message,
	redChan chan messagetargets.MessageTarget) interact.ContextInteract {
	i.db = database.GetDebugDB()

	source := msg.Source
	var user User = User{QQId: source.GetSenderID()}

	i.db.Where(&user).FirstOrCreate(&user)
	user.Records = append(user.Records, loadPayments(user, i.db)...)
	i.user = user

	i.status = Nil
	redChan <- messagetargets.SourceTarget(source, structs.NewTextChain(nilMsgText))
	return i
}

//NextMessage 向上下文提交信息
func (i *PaymentRecordInteract) NextMessage(
	msg structs.Message,
	redChan chan messagetargets.MessageTarget) {
	var chain []structs.MessageChainInfo
	if i.status == Nil {
		chain = i.onNil(msg)
	} else if i.status == ADD {
		chain = i.onAdd(msg)
	}
	if len(chain) > 0 {
		if msg.Source.GetSource() == constdata.GroupMessage {
			chain = append([]structs.MessageChainInfo{structs.NewAtChain(i.user.QQId, "")}, chain...)
		}

		redChan <- messagetargets.SourceTarget(msg.Source, chain...)
	}
	//finish operate
	//check status
	if i.status == Nil {
		redChan <- messagetargets.SourceTarget(msg.Source,
			structs.NewTextChain(nilMsgText),
		)
	}

}

//IsDone 判断该上下文是否已经完成
func (i *PaymentRecordInteract) IsDone() bool {
	return i.status == EXIT
}
