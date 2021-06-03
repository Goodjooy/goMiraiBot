package interact

import (
	messagetargets "goMiraiQQBot/lib/messageHandle/messageTargets"
	"goMiraiQQBot/lib/messageHandle/structs"
	"log"
)

/// interactHandler 管理全部的交互事件
type interactHandler struct {
	//全局上下文控制器
	activityContext ContextFetchMap

	//优先级列表
	prioityOrder []int32
	//管理具体控制器
	interactControllers map[int32]InteractController

	lock bool

	//信息流
	inputMsg chan structs.Message
	outputMsg chan messagetargets.MessageTarget
}
func (handle *interactHandler) DoLoadCmd(msg structs.Message) (
	Command,
	InteractController) {
	var cmd Command
	var controller InteractController
	for _, v := range handle.prioityOrder {
		controller, ok := handle.interactControllers[v]
		if ok {
			cmd, ok = controller.DoAnalyse(msg)
			if ok {
				return cmd, controller
			}
		}
	}
	return cmd, controller
}
func (handle *interactHandler) DoInteract(
	cmd Command,
	msg structs.Message,
	controller InteractController,
) bool {
	if interact, ok := controller.GetContextInteract(cmd, msg.Source); ok {
		context := interact.InitMessage(cmd.GetExtraCmd(), msg, handle.outputMsg)

		err := handle.activityContext.
			Put(msg.Source.GetGroupID(),
				msg.Source.GetSenderID(),
				context)

		if err != nil {
			log.Printf("Add New Context Failure : %v", err)
			handle.outputMsg <- messagetargets.SourceTarget(
				msg.Source,
				structs.NewTextChain("新建上下文失败"),
			)
		}
	} else if interact, ok := controller.GetSingleInteract(cmd, msg.Source); ok {
		interact.EnterMessage(
			cmd.GetExtraCmd(),
			msg,
			handle.outputMsg,
		)
	} else {
		return false
	}
	return true
}

func (handle *interactHandler) DoContextAlive(
	msg structs.Message,
) bool {
	if context, err := handle.activityContext.Get(msg.Source.GetGroupID(), msg.Source.GetSenderID()); err == nil {
		context.NextMessage(msg, handle.outputMsg)

		if context.IsDone() {
			handle.activityContext.Delete(msg.Source.GetGroupID(), msg.Source.GetSenderID())
		}
		return true
	}
	return false
}
