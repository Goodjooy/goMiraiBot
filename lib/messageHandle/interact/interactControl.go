package interact

import (
	"goMiraiQQBot/lib/messageHandle/sourceHandle"
	"goMiraiQQBot/lib/messageHandle/structs"
)

type InteractController interface {
	//GetPrioity 获取当前交互控制器的优先级。
	//优先级相同的控制器不能保证严格先后顺序
	GetPrioity() int32

	//DoAnalyse 根据输入的Message 信息，分析指令。
	//如果找不到指令，返回false
	DoAnalyse(structs.Message) (Command, bool)

	ConstructorContain
}

type ConstructorContain interface {
	//GetSingleInteract 通过 Cmd,MessageSource 获取指定的单次交互interact
	//如果没有找到，返回false
	GetSingleInteract(Command, sourceHandle.MessageSource) (FullSingleInteract, bool)
	//GetContextInteract 通过 Cmd,MessageSource 获取指定的上下文交互interact
	//如果没有找到，返回false
	GetContextInteract(Command, sourceHandle.MessageSource) (FullContextInteract, bool)

	//AddSingleInteractConstruct 添加interact构造器，向interact存储容器里面添加
	AddSingleInteractConstruct(SingleInteractConstruct)
	//AddContextInteractConstruct 添加interact构造器，向interact存储容器里面添加
	AddContextInteractConstruct(ContextInteractConstruct)
	//LockAdder 锁定构造器添加
	LockAdder()
}

type PerfabInteractContorller struct {
	prioity int32

	analyseFn func(structs.Message) (Command, bool)

	PerfabContructorContainer
}

func NewInteractContorller(prioity int32, analyseFn func(structs.Message) (Command, bool)) InteractController {
	return &PerfabInteractContorller{
		prioity:                   prioity,
		analyseFn:                 analyseFn,
		PerfabContructorContainer: newPerfabConstructorContainer(),
	}
}

//GetPrioity 获取当前交互控制器的优先级。
//优先级相同的控制器不能保证严格先后顺序
func (i *PerfabInteractContorller) GetPrioity() int32 {
	return i.prioity
}

//DoAnalyse 根据输入的Message 信息，分析指令。
//如果找不到指令，返回false
func (i *PerfabInteractContorller) DoAnalyse(msg structs.Message) (Command, bool) {
	return i.analyseFn(msg)
}

type PerfabContructorContainer struct {
	single  ConstructMap
	context ConstructMap

	lock bool
}

func newPerfabConstructorContainer() PerfabContructorContainer {
	return PerfabContructorContainer{
		single:  NewContructMap(),
		context: NewContructMap(),
		lock:    false,
	}
}

//GetSingleInteract 通过 Cmd,MessageSource 获取指定的单次交互interact
//如果没有找到，返回false
func (p *PerfabContructorContainer) GetSingleInteract(cmd Command, source sourceHandle.MessageSource) (FullSingleInteract, bool) {
	interact, err := p.single.GetSingleInteract(cmd.mainCmd, source.GetSource())
	if err != nil {
		return nil, false
	}
	return interact, true
}

//GetContextInteract 通过 Cmd,MessageSource 获取指定的上下文交互interact
//如果没有找到，返回false
func (p *PerfabContructorContainer) GetContextInteract(cmd Command, source sourceHandle.MessageSource) (FullContextInteract, bool) {
	interact, err := p.context.GetContextInteract(cmd.mainCmd, source.GetSource())
	if err != nil {
		return nil, false
	}
	return interact, true
}

//AddSingleInteractConstruct 添加interact构造器，向interact存储容器里面添加
func (p *PerfabContructorContainer) AddSingleInteractConstruct(interact SingleInteractConstruct) {
	if p.lock {
		panic("The Container has been locked")
	}
	p.single.AddSingleConstruct(interact)
}

//AddContextInteractConstruct 添加interact构造器，向interact存储容器里面添加
func (p *PerfabContructorContainer) AddContextInteractConstruct(interact ContextInteractConstruct) {
	if p.lock {
		panic("The Container has been locked")
	}
	p.context.AddContextConstruct(interact)
}

//LockAdder 锁定构造器添加
func (p *PerfabContructorContainer) LockAdder() {
	p.lock = true
	p.single.setLock()
	p.context.setLock()
}
