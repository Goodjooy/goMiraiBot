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

