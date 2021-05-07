package interact

import (
	"goMiraiQQBot/constdata"
	cmi "goMiraiQQBot/messageHandle/interact/contextMessageInteract"
	smi "goMiraiQQBot/messageHandle/interact/sinpleMessageInteract"
	msh "goMiraiQQBot/messageHandle/messageSourceHandles"
	messagetargets "goMiraiQQBot/messageHandle/messageTargets"
	"goMiraiQQBot/messageHandle/structs"
	"regexp"
	"strings"
)

type GroupMemberContext map[uint64]cmi.ContextMessageInteract

type SingleInteractConstruct func() smi.SingleMessageInteract
type ContextInteractConstruct func() cmi.ContextMessageInteract

//活跃的信息注册
var activateContextInteract map[uint64]GroupMemberContext = make(map[uint64]GroupMemberContext)

//构造容器
var singleInteract map[string]SingleInteractConstruct = make(map[string]SingleInteractConstruct)
var contextInteract map[string]ContextInteractConstruct = map[string]ContextInteractConstruct{}

func GetSingleCommand()[]string{
	var cmds []string
	for k := range singleInteract {cmds = append(cmds, k)
	}
	return cmds
}
func GetContextCommand()[]string{
	var cmds []string
	for k := range contextInteract {cmds = append(cmds, k)
	}
	return cmds
}

func InitInteractHandle(msgChan chan structs.Message, msgRes chan messagetargets.MessageTarget) {
	//register interacter
	addSingleInteract(NewHelpInteract)
	addSingleInteract(smi.NewAboutInteract)

	
	go acceptMessage(msgChan, msgRes)
}

func addSingleInteract(handle SingleInteractConstruct){
	key:=strings.ToLower(handle().GetCommandName())

	singleInteract[key]=handle
}

func acceptMessage(msgChan chan structs.Message, msgRes chan messagetargets.MessageTarget) {
	for {
		select {
		case data, ok := (<-msgChan):
			if ok {
				source := data.Source
				//信息类型为群消息
				if source.GetSource() == constdata.GroupMessage {
					var d msh.GroupMessage = source.GetMetaInformation().(msh.GroupMessage)

					//上下文环境持续，最高优先级
					if group, ok := activateContextInteract[d.GroupId]; ok {
						if context, ok := group[d.UserId]; ok {
							context.NextMessage(data, msgRes)
							//check done
							if context.IsDone() {
								delete(group, d.UserId)
							}
						}
					} else {
						//激活上下文，第二优先级
						msgChain := data.ChainInfoList
						cmd, ok := commandGet(msgChain)
						if ok {
							if context, ok := contextInteract[cmd]; ok {
								var c = context().InitMessage(
									data,
									msgRes,
								)

								group, ok := activateContextInteract[d.GroupId]
								if !ok {
									group := make(GroupMemberContext)
									activateContextInteract[d.GroupId] = group
								}
								group[d.UserId] = c
							} else if signle, ok := singleInteract[cmd]; ok {
								signle().EnterMessage(data, msgRes)
							}
						}
					}
				}
			}
		}
	}
}

var cmdPattern = regexp.MustCompile(`^#\s*(\S+)\s*`)

func commandGet(msgChain []structs.MessageChainInfo) (string, bool) {
	for _, v := range msgChain {
		//找到第一段文本
		if v.MessageType == constdata.Plain {
			//获取文本
			msg := v.Data["text"].(string)
			//以#开头
			if strings.HasPrefix(msg, "#") {
				cmd := cmdPattern.FindStringSubmatch(msg)[1]
				return strings.ToLower(cmd), true
			}
		}
	}
	return "", false
}
