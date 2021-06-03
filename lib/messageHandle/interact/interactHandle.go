package interact

import (
	messagetargets "goMiraiQQBot/lib/messageHandle/messageTargets"
	"goMiraiQQBot/lib/messageHandle/sourceHandle"
	"goMiraiQQBot/lib/messageHandle/structs"
	"log"
)

type BotQQIdGeter interface {
	GetQQId() uint64
}

func SetCFG(c BotQQIdGeter){
	cfg=c
}

func InitInteractHandle(msgChan chan structs.Message, msgRes chan messagetargets.MessageTarget) {

	handler=interactHandler{
		activityContext: NewContextFetchMap(),
		
	}

	MessageInteract.setLock()
	ChainInteract.setLock()
	//register interacter
	go acceptMessage(msgChan, msgRes)
	go acceptMessage(msgChan, msgRes)
}
func acceptMessage(msgChan chan structs.Message, msgRes chan messagetargets.MessageTarget) {
	for {
		select {
		case data, ok := (<-msgChan):
			if ok {
				source := data.Source

				//上下文环境持续，最高优先级
				status := contiuneContextHandle(source, data, msgRes)
				if status {
					continue
				}
				// 命令形式上下文
				msgChain := data.ChainInfoList
				cmd, ok := CommandGet(msgChain, cfg.GetQQId())
				if ok {
					interactActivateHandle(cmd, source, data, msgRes, true, MessageInteract)
					continue
				}

				//类型形式上下文
				chainCmd, ok := chainStructGet(msgChain)
				if ok {
					interactActivateHandle(chainCmd, source, data, msgRes, false, ChainInteract)
					continue
				}
			}
		}
	}
}

func interactActivateHandle(cmd Command,
	source sourceHandle.MessageSource,
	data structs.Message,
	msgRes chan messagetargets.MessageTarget,
	cmdNfoundMsg bool,
	constructContainer ConstructMap) {

	if context, err := constructContainer.GetContextInteract(cmd.mainCmd, source.GetSource()); err == nil {
		var c = context.InitMessage(
			cmd.extraCmd,
			data,
			msgRes,
		)
		err := activateContextInteract.Put(source.GetGroupID(), source.GetSenderID(), c)
		if err != nil {
			log.Printf("Add New Context Fauilure: %v", err)
			msgRes <- messagetargets.NewSingleTextGroupTarget(source.GetSenderID(), "新建上下文失败")
		}
	} else if signle, err := constructContainer.GetSingleInteract(cmd.mainCmd, source.GetSource()); err == nil {
		signle.EnterMessage(cmd.extraCmd, data, msgRes)
	} else {
		if cmdNfoundMsg {
			msgRes <- messagetargets.NewSingleTextGroupTarget(source.GetGroupID(), "指令未找到！")
		}
	}
}

func contiuneContextHandle(
	source sourceHandle.MessageSource,
	data structs.Message,
	msgRes chan messagetargets.MessageTarget,
) bool {
	if context, err := activateContextInteract.Get(source.GetGroupID(), source.GetSenderID()); err == nil {
		context.NextMessage(data, msgRes)
		//check done
		if context.IsDone() {
			activateContextInteract.Delete(source.GetGroupID(), source.GetSenderID())
		}
		return true
	}
	return false
}
