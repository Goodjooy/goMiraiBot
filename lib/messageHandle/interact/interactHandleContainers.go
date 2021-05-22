package interact

import (
	"fmt"
	"sync"

	uuid "github.com/satori/go.uuid"
)

type GroupMemberContext map[uint64]ContextInteract



type ContextFetchMap struct {
	//数据
	contextData     map[uuid.UUID]ContextInteract
	groupMenberData map[uint64]map[uint64]uuid.UUID

	//互斥锁
	mutex sync.RWMutex
}
func NewContextFetchMap()ContextFetchMap{
	return ContextFetchMap{
		contextData: make(map[uuid.UUID]ContextInteract),
		groupMenberData: make(map[uint64]map[uint64]uuid.UUID),
		mutex: sync.RWMutex{},
	}
}

//Put put add new ContextInteract into target group and menber area, <br>
//if target area has not done Context will returen an error  <br>
//while adding Context will lock by mutex  <br>
func (c *ContextFetchMap) Put(groupId, menberId uint64, context ContextInteract) error {
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
func (c *ContextFetchMap) Get(groupId, menberId uint64) (ContextInteract, error) {

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
func (c *ContextFetchMap) Delete(groupId, menberId uint64) error {
	contextUUID, err := c.getContextUUID(groupId, menberId)
	if err != nil {
		return fmt.Errorf("Get Context UUID Error: %v", err)
	}

	c.mutex.Lock()
	delete(c.contextData, contextUUID)
	c.mutex.Unlock()
	return nil
}

func (c *ContextFetchMap) getContextUUID(groupId, menberId uint64) (uuid.UUID, error) {
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
var activateContextInteract ContextFetchMap = NewContextFetchMap()

//构造容器
var MessageInteract ConstructMap=NewContructMap()

var ChainInteract ConstructMap=NewContructMap()

var cfg BotQQIdGeter

func GetBotQQ()uint64{
	return cfg.GetQQId()
}