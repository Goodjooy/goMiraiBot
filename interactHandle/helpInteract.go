package interactHandle

import (
	"goMiraiQQBot/constdata"
	datautil "goMiraiQQBot/dataUtil"
	"goMiraiQQBot/messageHandle/interact"
	messagetargets "goMiraiQQBot/messageHandle/messageTargets"
	"goMiraiQQBot/messageHandle/sourceHandle"
	"goMiraiQQBot/messageHandle/structs"
	"strings"

	uuid "github.com/satori/go.uuid"
)

var (
	command    = datautil.NewTargetValues("command", "指令", "命令", "cmd", "指令型")
	typeActive = datautil.NewTargetValues("type", "信息类型", "类型")
)

var (
	cmd    = datautil.NewTargetValues("name", "指令", "命令")
	source = datautil.NewTargetValues("type", "来源", "类型")
)

//帮助功能，用于显示全部命令
type HelpInteract struct {
}

func NewHelpInteract() interact.FullSingleInteract {
	return HelpInteract{}
}
func (HelpInteract) Init() {

}

func (interact HelpInteract) GetCommandName() []string {
	return []string{"help", "帮助", "功能"}
}
func (interact HelpInteract) RespondSource() []constdata.MessageType {
	return []constdata.MessageType{
		constdata.GroupMessage,
		constdata.FriendMessage,
	}
}

func (i HelpInteract) EnterMessage(
	extraCmd datautil.MutliToOneMap,
	data structs.Message,
	repChan chan messagetargets.MessageTarget) {

	extraCmd.SetNoNameCmdOrder(cmd, source)
	var target messagetargets.MessageTarget
	if extraCmd.IsEmpty() {
		target = getAllCmd(data.Source)
	} else {
		cmdName, nok := extraCmd.Get(cmd...)
		cmdName = strings.ToLower(cmdName)
		sourceName, _ := extraCmd.GetWithDefault(command.GetSign(), source...)

		if nok {
			target = getTargetCmd(sourceName, cmdName, data.Source)
		} else {
			rep := messagetargets.SourceTarget(data.Source, structs.NewTextChain("指令不完整！name 和 type 要同时指定"))
			target = rep
		}
	}
	repChan <- target
}

func (HelpInteract) GetUseage() string {
	return "#help|#帮助 : 返回当前机器人的全部指令\n" +
		"额外指令：\n" +
		" name|指令|命令=[指令名称] \n  -> 标记选择的指令名称\n" +
		" type|来源|类型=[single|context] \n  -> 选择指令来源【单次交互|上下文交互】"
}

func cmdString(prefix, sep string) []structs.MessageChainInfo {
	var d []structs.MessageChainInfo

	var lastCmdUUID uuid.UUID
	var start bool = false

	for _, v := range interact.MessageInteract.GetAllCmdSet() {
		if start && lastCmdUUID == v.GetUUID() {
			continue
		}
		var supportCmd string

		supportCmd += (prefix + " ")
		for _, c := range v.GetCmds() {
			supportCmd += (c + " " + sep + " ")
		}
		supportCmd = supportCmd[:len(supportCmd)-len(sep)]

		supportCmd += "\n"
		d = append(d, structs.NewTextChain(supportCmd))

		start = true
		lastCmdUUID = v.GetUUID()
	}
	return d
}

func getAllCmd(source sourceHandle.MessageSource) messagetargets.MessageTarget {
	var d []structs.MessageChainInfo

	d = append(d, structs.NewTextChain("指令交互命令：\n"))

	d = append(d, cmdString("#", "|")...)

	d = append(d, structs.NewTextChain("信息类型响应\n"))

	d = append(d, cmdString("", "|")...)

	da := messagetargets.SourceTarget(source, d...)

	return da
}

func getTargetCmd(sourceName, cmdName string, source sourceHandle.MessageSource) messagetargets.MessageTarget {
	if command.Match(sourceName) {
		c, ok := interact.MessageInteract.GetSideInfoFromCmd(cmdName)
		if ok {
			rep := messagetargets.SourceTarget(source,
				structs.NewTextChain(c.GetUseage()))
			return rep
		}
	} else if typeActive.Match(sourceName) {
		c, ok := interact.ChainInteract.GetSideInfoFromCmd(cmdName)
		if ok {
			rep := messagetargets.SourceTarget(source, structs.NewTextChain(c.GetUseage()))
			return rep
		}
	}
	rep := messagetargets.SourceTarget(source, structs.NewTextChain("指令未找到！"))
	return rep
}
