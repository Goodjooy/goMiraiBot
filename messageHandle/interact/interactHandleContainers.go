package interact

import (
	messagetargets "goMiraiQQBot/messageHandle/messageTargets"
	"goMiraiQQBot/messageHandle/structs"
	"strings"
	"sync"
)

type GroupMemberContext map[uint64]ContextInteract

type SingleInteractConstruct func() SingleMessageInteract
type ContextInteractConstruct func() ContextMessageInteract

type ChainSingleInteractConstruct func() ChainTypeInteract
type ChainContextInteractConstruct func() ChainTypeContextInteract

//活跃的信息注册
var activateContextInteract map[uint64]GroupMemberContext = make(map[uint64]GroupMemberContext)
var activateMut sync.RWMutex

//构造容器
var singleInteract map[string]SingleInteractConstruct = make(map[string]SingleInteractConstruct)
var contextInteract map[string]ContextInteractConstruct = make(map[string]ContextInteractConstruct)

var chainSingleInteact map[string]ChainSingleInteractConstruct = make(map[string]ChainSingleInteractConstruct)
var chainContextInteract map[string]ChainContextInteractConstruct = make(map[string]ChainContextInteractConstruct)

var cfg BotQQIdGeter

func GetSingleCommand() []string {
	var cmds []string
	for k := range singleInteract {
		cmds = append(cmds, k)
	}
	return cmds
}
func GetContextCommand() []string {
	var cmds []string
	for k := range contextInteract {
		cmds = append(cmds, k)
	}
	return cmds
}
func GetChainContextCommand() []string {
	var cmds []string
	for k := range chainContextInteract {
		cmds = append(cmds, k)
	}
	return cmds
}
func GetChainSingleCommand() []string {
	var cmds []string
	for k := range chainSingleInteact {
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

func InitInteractHandle(msgChan chan structs.Message, msgRes chan messagetargets.MessageTarget, c BotQQIdGeter) {
	cfg = c

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

func AddChainSingleInteract(handle ChainSingleInteractConstruct) {
	keys := handle().GetActivateType()

	for _, key := range keys {
		k := strings.ToLower(key.String())
		chainSingleInteact[k] = handle
	}

}
func AddChainContextInteract(handle ChainContextInteractConstruct) {
	keys := handle().GetActivateType()

	for _, key := range keys {
		k := strings.ToLower(key.String())
		chainContextInteract[k] = handle
	}
}
