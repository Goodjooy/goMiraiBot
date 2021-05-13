package interact

import (
	"goMiraiQQBot/constdata"
	messagetargets "goMiraiQQBot/messageHandle/messageTargets"
	"goMiraiQQBot/messageHandle/sourceHandle"
	"goMiraiQQBot/messageHandle/structs"
	"log"
)

func acceptMessage(msgChan chan structs.Message, msgRes chan messagetargets.MessageTarget) {
	for {
		select {
		case data, ok := (<-msgChan):
			if ok {
				source := data.Source
				//信息类型为群消息
				if source.GetSource() == constdata.GroupMessage {
					var d sourceHandle.GroupMessage = source.GetMetaInformation().(sourceHandle.GroupMessage)

					//上下文环境持续，最高优先级
					if context, err := activateContextInteract.Get(d.GroupId, d.UserId); err == nil {
						context.NextMessage(data, msgRes)
						//check done
						if context.IsDone() {
							activateContextInteract.Delete(d.GroupId, d.UserId)
						}
					}
					//激活上下文，第二优先级
					msgChain := data.ChainInfoList
					cmd, ok := commandGet(msgChain)
					if ok {
						if context, ok := contextInteract[cmd.mainCmd]; ok {
							var c = context().InitMessage(
								cmd.extraCmd,
								data,
								msgRes,
							)
							err := activateContextInteract.Put(d.GroupId, d.UserId, c)
							if err != nil {
								log.Printf("Add New Context Faulre: %v", err)
								//TODO: error Channal
								msgRes <- messagetargets.NewSingleTextGroupTarget(d.GroupId, "新建上下文失败")
							}
						} else if signle, ok := singleInteract[cmd.mainCmd]; ok {
							signle().EnterMessage(cmd.extraCmd, data, msgRes)
							continue
						} else {
							msgRes <- messagetargets.NewSingleTextGroupTarget(d.GroupId, "指令未找到！")
						}
					}
					chainCmd, ok := chainStructGet(msgChain)
					if ok {
						if context, ok := chainContextInteract[chainCmd.mainCmd]; ok {
							var c = context().InitMessage(
								cmd.extraCmd,
								data,
								msgRes,
							)
							err := activateContextInteract.Put(d.GroupId, d.UserId, c)
							if err != nil {
								log.Printf("Add New Context Faulre: %v", err)
								//TODO: error Channal
								msgRes <- messagetargets.NewSingleTextGroupTarget(d.GroupId, "新建上下文失败")
							}
						} else if signle, ok := chainSingleInteact[chainCmd.mainCmd]; ok {
							signle().EnterMessage(cmd.extraCmd, data, msgRes)
						}
					}
				}
			}
		}
	}
}
