package randomsetuinteract

import (
	"goMiraiQQBot/lib/constdata"
	datautil "goMiraiQQBot/lib/dataUtil"
	interactprefab "goMiraiQQBot/lib/interactPrefab"
	"goMiraiQQBot/lib/messageHandle/interact"
	messagetargets "goMiraiQQBot/lib/messageHandle/messageTargets"
	"goMiraiQQBot/lib/messageHandle/structs"
)


type SetuInteract struct {
	interactprefab.InteractPerfab
}

func NewSetuInteract()interact.FullSingleInteract{
	perfab:=interactprefab.
	NewInteractPerfab().

	AddActivateSigns("setu").
	AddActivateSigns("随机涩图").
	AddActivateSigns("涩图").

	AddActivateSource(constdata.GroupMessage).
	AddActivateSource(constdata.FriendMessage).
	
	AddInitFn(func() {}).
	
	SetUseage(`随机涩图功能，每日额度300`).

	BuildPtr()

	return &SetuInteract{InteractPerfab: *perfab}
}

func(setu*SetuInteract)EnterMessage(
	extraCmd datautil.MutliToOneMap, 
	data structs.Message, 
	repChan chan messagetargets.MessageTarget){
	setu.TodoMsgRespond(data.Source,repChan)
}