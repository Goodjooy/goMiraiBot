package inandoutinteract

import (
	"goMiraiQQBot/lib/constdata"
	datautil "goMiraiQQBot/lib/dataUtil"
	"goMiraiQQBot/lib/database"
	interactprefab "goMiraiQQBot/lib/interactPrefab"
	"goMiraiQQBot/lib/messageHandle/interact"
	messagetargets "goMiraiQQBot/lib/messageHandle/messageTargets"
	"goMiraiQQBot/lib/messageHandle/structs"

	"gorm.io/gorm"
)

var botQQ uint64

type PaymentRecordInteract struct {
	interactprefab.InteractPerfab

	user User
	db   *gorm.DB

	status Status
}

func NewPaymentRecordInteract() interact.FullContextInteract {

	return &PaymentRecordInteract{
		InteractPerfab: interactprefab.NewInteractPerfab().
			AddActivateSigns("记账", "账单").
			AddActivateSigns("payment-record").
			SetUseage(`#记账 - 上下文类型命令交互
		在之后可进行多种指令操作`).
			AddActivateSource(constdata.GroupMessage, constdata.FriendMessage).
			AddInitFn(func() {
				database.AsignDBModel(&PaymentRecord{}, &User{})
				botQQ = interact.GetBotQQ()
			}).
			Build(),
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
