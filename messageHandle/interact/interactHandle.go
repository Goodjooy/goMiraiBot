package interact

import (
	"goMiraiQQBot/constdata"
	messagetargets "goMiraiQQBot/messageHandle/messageTargets"
	"goMiraiQQBot/messageHandle/sourceHandle"
	"goMiraiQQBot/messageHandle/structs"
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
					if group, ok := activateContextInteract[d.GroupId]; ok {
						if context, ok := group[d.UserId]; ok {
							context.NextMessage(data, msgRes)
							//check done
							if context.IsDone() {
								activateMut.Lock()
								delete(group, d.UserId)
								activateMut.Unlock()
							}
							continue
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

							group, ok := activateContextInteract[d.GroupId]
							if !ok {
								group = make(GroupMemberContext)
								activateMut.Lock()
								activateContextInteract[d.GroupId] = group
								activateMut.Unlock()
							}
							activateMut.Lock()
							group[d.UserId] = c
							activateMut.Unlock()
							continue
						} else if signle, ok := singleInteract[cmd.mainCmd]; ok {
							signle().EnterMessage(cmd.extraCmd, data, msgRes)
							continue
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
							group, ok := activateContextInteract[d.GroupId]
							if !ok {
								group = make(GroupMemberContext)
								activateMut.Lock()
								activateContextInteract[d.GroupId] = group
								activateMut.Unlock()
							}
							activateMut.Lock()
							group[d.UserId] = c
							activateMut.Unlock()
						} else if signle, ok := chainSingleInteact[chainCmd.mainCmd]; ok {
							signle().EnterMessage(cmd.extraCmd, data, msgRes)
						}
					}
				}
			}
		}
	}
}
