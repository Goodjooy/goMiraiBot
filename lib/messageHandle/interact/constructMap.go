package interact

import (
	"errors"
	"fmt"
	"goMiraiQQBot/lib/constdata"
	"log"

	uuid "github.com/satori/go.uuid"
)

type SingleInteractConstruct func() FullSingleInteract
type ContextInteractConstruct func() FullContextInteract
type InteractSideConstruct func() InteractSideInformation

type InteractSideInformation interface {
	//GetUseage
	GetUseage() string

	// GetCommandName 获取指令名称，用来创建映射关系
	GetCommandName() []string
	//RespondSource
	RespondSource() []constdata.MessageType
}

type CmdSet struct {
	cmds     []string
	msgTypes []constdata.MessageType

	uid uuid.UUID
}

func NewCmdSet(msgTypes []constdata.MessageType, cmmands []string) CmdSet {
	return CmdSet{
		cmds:     cmmands,
		msgTypes: msgTypes,
		uid:      uuid.NewV1(),
	}
}

func (c CmdSet) GetCmds() []string {
	return c.cmds
}
func (c CmdSet) GetMsgType() []constdata.MessageType {
	return c.msgTypes
}
func (c CmdSet) GetUUID() uuid.UUID {
	return c.uid
}

type ConstructMap struct {
	sideInfos map[uuid.UUID]InteractSideConstruct

	sinpleConstruct  map[uuid.UUID]SingleInteractConstruct
	contextConstruct map[uuid.UUID]ContextInteractConstruct

	cmdMap map[string]CmdSet

	lock bool
}

func NewContructMap() ConstructMap {
	return ConstructMap{
		sideInfos:        make(map[uuid.UUID]InteractSideConstruct),
		sinpleConstruct:  make(map[uuid.UUID]SingleInteractConstruct),
		contextConstruct: make(map[uuid.UUID]ContextInteractConstruct),
		cmdMap:           make(map[string]CmdSet),
		lock:             false,
	}
}
func (c *ConstructMap) AddSingleConstruct(construct SingleInteractConstruct) {
	if c.lock {
		log.Panicf("add Construct After Locked!")
	}
	interact := construct()
	interact.Init()
	cmdSet := c.AddCmdSet(interact)

	c.sideInfos[cmdSet.uid] = func() InteractSideInformation { return construct() }
	c.sinpleConstruct[cmdSet.uid] = construct

}

func (c *ConstructMap) AddContextConstruct(construct ContextInteractConstruct) {
	if c.lock {
		log.Panicf("add Construct After Locked!")
	}
	interact := construct()
	interact.Init()
	cmdSet := c.AddCmdSet(interact)

	c.sideInfos[cmdSet.uid] = func() InteractSideInformation { return construct() }
	c.contextConstruct[cmdSet.uid] = construct
}

func (c *ConstructMap) AddCmdSet(sideInfo InteractSideInformation) CmdSet {
	if c.lock {
		log.Panicf("add Construct After Locked!")
	}
	cmdSet := NewCmdSet(sideInfo.RespondSource(), sideInfo.GetCommandName())
	for _, v := range sideInfo.GetCommandName() {
		c.cmdMap[v] = cmdSet
	}
	return cmdSet
}

func (c *ConstructMap) GetSingleInteract(cmd string, msgType constdata.MessageType) (FullSingleInteract, error) {
	cmdSet, err := c.GetCmdSet(cmd, msgType)
	if err != nil {
		return nil, err
	}
	construct, ok := c.sinpleConstruct[cmdSet.uid]
	if ok {
		return construct(), nil
	}
	return nil, errors.New("no Found Target Construct")
}

func (c *ConstructMap) GetContextInteract(cmd string, msgType constdata.MessageType) (FullContextInteract, error) {
	cmdSet, err := c.GetCmdSet(cmd, msgType)
	if err != nil {
		return nil, err
	}
	construct, ok := c.contextConstruct[cmdSet.uid]
	if ok {
		return construct(), nil
	}
	return nil, errors.New("no Found Target Construct")
}

func (c *ConstructMap) GetCmdSetWithOutSourceLimit(cmd string) (CmdSet, error) {
	cmdSet, ok := c.cmdMap[cmd]
	if !ok {
		return CmdSet{}, fmt.Errorf("no Activate Sign Found For `%s`", cmd)
	}
	return cmdSet, nil
}

func (c *ConstructMap) GetCmdSet(cmd string, msgType constdata.MessageType) (CmdSet, error) {
	cmdSet, err := c.GetCmdSetWithOutSourceLimit(cmd)
	if err != nil {
		return CmdSet{}, err
	}
	//search for msgType
	for _, v := range cmdSet.msgTypes {
		if v == msgType {
			goto SUPPORT
		}
	}
	return CmdSet{}, fmt.Errorf("message Sourece<%s> Not Support Target Activate Sign", msgType)

SUPPORT:
	{
		return cmdSet, nil
	}
}

func (c *ConstructMap) GetAllCmdSet() []CmdSet {
	var cmdSets []CmdSet
	for _, v := range c.cmdMap {
		cmdSets = append(cmdSets, v)
	}
	return cmdSets
}
func (c *ConstructMap) GetSideInfoFromCmd(cmd string) (InteractSideInformation, bool) {
	cmdSet, err := c.GetCmdSetWithOutSourceLimit(cmd)
	if err != nil {
		log.Print("Get Side Info Failure: ",err)
		return nil, false
	}
	sideInfo, err := c.GetSideInfo(cmdSet)
	if err != nil {
		log.Print("Get Side Info Failure: ",err)
		return nil, false
	}
	return sideInfo, true
}

func (c *ConstructMap) GetSideInfo(cmdSet CmdSet) (InteractSideInformation, error) {
	sideInfoConstruct, ok := c.sideInfos[cmdSet.uid]
	if !ok {
		return nil, fmt.Errorf("no Side Info Construct For %v", cmdSet.cmds[0])
	}
	return sideInfoConstruct(), nil
}

func (c *ConstructMap) setLock() {
	c.lock = true
}
