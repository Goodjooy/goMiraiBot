package interact

import (
	"goMiraiQQBot/constdata"
	messagetargets "goMiraiQQBot/messageHandle/messageTargets"
	"goMiraiQQBot/messageHandle/sourceHandle"
	"goMiraiQQBot/messageHandle/structs"
	"strings"
	"sync"
)

type GroupMemberContext map[uint64]ContextMessageInteract

type SingleInteractConstruct func() SingleMessageInteract
type ContextInteractConstruct func() ContextMessageInteract

//活跃的信息注册
var activateContextInteract map[uint64]GroupMemberContext = make(map[uint64]GroupMemberContext)
var activateMut sync.RWMutex

//构造容器
var singleInteract map[string]SingleInteractConstruct = make(map[string]SingleInteractConstruct)
var contextInteract map[string]ContextInteractConstruct = map[string]ContextInteractConstruct{}

func SetSingleCommand() []string {
	var cmds []string
	for k := range singleInteract {
		cmds = append(cmds, k)
	}
	return cmds
}
func SetContextCommand() []string {
	var cmds []string
	for k := range contextInteract {
		cmds = append(cmds, k)
	}
	return cmds
}

func GetSingleInteract(key string) (SingleInteractConstruct, bool) {
	v, ok := singleInteract[key]
	return v, ok
}

func GetContextInteract(key string) (ContextInteractConstruct, bool) {
	v, ok := contextInteract[key]
	return v, ok
}

func InitInteractHandle(msgChan chan structs.Message, msgRes chan messagetargets.MessageTarget) {
	//register interacter
	go acceptMessage(msgChan, msgRes)
	go acceptMessage(msgChan, msgRes)
}

func AddSingleInteract(handle SingleInteractConstruct) {
	keys := handle().GetCommandName()

	for _, key := range keys {
		key = strings.ToLower(key)
		singleInteract[key] = handle
	}

}
func AddContextInteract(handle ContextInteractConstruct) {
	keys := handle().GetInitCommand()

	for _, key := range keys {
		key = strings.ToLower(key)
		contextInteract[key] = handle
	}
}

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
						} else if signle, ok := singleInteract[cmd.mainCmd]; ok {
							signle().EnterMessage(cmd.extraCmd, data, msgRes)
						}
					}
				}
			}
		}
	}
}
