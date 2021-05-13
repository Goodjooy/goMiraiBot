package interact

import (
	"fmt"
	messagetargets "goMiraiQQBot/messageHandle/messageTargets"
	"goMiraiQQBot/messageHandle/structs"
	"strings"
	"sync"

	uuid "github.com/satori/go.uuid"
)

type GroupMemberContext map[uint64]ContextInteract

type SingleInteractConstruct func() SingleMessageInteract
type ContextInteractConstruct func() ContextMessageInteract

type ChainSingleInteractConstruct func() ChainTypeInteract
type ChainContextInteractConstruct func() ChainTypeContextInteract

type contextFetchMap struct {
	//数据
	contextData     map[uuid.UUID]ContextInteract
	groupMenberData map[uint64]map[uint64]uuid.UUID

	//互斥锁
	mutex sync.RWMutex
}
func newContextFetchMap()contextFetchMap{
	return contextFetchMap{
		contextData: make(map[uuid.UUID]ContextInteract),
		groupMenberData: make(map[uint64]map[uint64]uuid.UUID),
		mutex: sync.RWMutex{},
	}
}

//Put put add new ContextInteract into target group and menber area, <br>
//if target area has not done Context will returen an error  <br>
//while adding Context will lock by mutex  <br>
func (c *contextFetchMap) Put(groupId, menberId uint64, context ContextInteract) error {
	oldContext, err := c.Get(groupId, menberId)
	if err == nil {
		if !oldContext.IsDone() {
			//last context not done
			//error
			return fmt.Errorf("last Context Not DONE In Group<ID:%v> For Member<ID:%v>", groupId, menberId)
		}
	}
	// context not exist
	defer c.mutex.Unlock()
	c.mutex.Lock()

	//save Context
	contextUUID := uuid.NewV1()
	c.contextData[contextUUID] = context

	menberMap, ok := c.groupMenberData[groupId]
	if !ok {
		//not exist group, create
		menberMap = make(map[uint64]uuid.UUID)
		c.groupMenberData[groupId] = menberMap
	}
	menberMap[menberId] = contextUUID

	return nil
}

//Get get get target context in provid group for the menber
//if context not found, will returen an error
//wihle get it will enable READ LOCK
func (c *contextFetchMap) Get(groupId, menberId uint64) (ContextInteract, error) {

	contextUUID, err := c.getContextUUID(groupId, menberId)
	if err != nil {
		return nil, err
	}

	c.mutex.RLock()
	//read context
	context, ok := c.contextData[contextUUID]
	c.mutex.RUnlock()

	if !ok {
		return nil, fmt.Errorf("no Context Found For User<ID:%v> Group<ID:%v>", menberId, groupId)
	}
	return context, nil
}

//Delete delete target Context no matter The Target Context is DONE or not
// while Delete, it will enable LOCK
func (c *contextFetchMap) Delete(groupId, menberId uint64) error {
	contextUUID, err := c.getContextUUID(groupId, menberId)
	if err != nil {
		return fmt.Errorf("Get Context UUID Error: %v", err)
	}

	c.mutex.Lock()
	delete(c.contextData, contextUUID)
	c.mutex.Unlock()
	return nil
}

func (c *contextFetchMap) getContextUUID(groupId, menberId uint64) (uuid.UUID, error) {
	defer c.mutex.Unlock()
	c.mutex.Lock()

	menberMap, ok := c.groupMenberData[groupId]
	if !ok {
		return uuid.UUID{}, fmt.Errorf("no Group Found For GroupID: %v", groupId)
	}
	contextUUID, ok := menberMap[menberId]
	if !ok {
		return uuid.UUID{}, fmt.Errorf("no Such User<ID:%v> Context In Group<ID:%v>", menberId, groupId)
	}
	return contextUUID, nil
}

//活跃的信息注册
var activateContextInteract contextFetchMap = newContextFetchMap()

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
