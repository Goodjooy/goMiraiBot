package interactHandle

import (
	"goMiraiQQBot/constdata"
	datautil "goMiraiQQBot/dataUtil"
	"goMiraiQQBot/messageHandle/interact"
	messagetargets "goMiraiQQBot/messageHandle/messageTargets"
	"goMiraiQQBot/messageHandle/sourceHandle"
	"goMiraiQQBot/messageHandle/structs"
	"goMiraiQQBot/request"
	"strings"
)

var (
	single     = datautil.NewTargetValues("single", "单次", "单次交互", "单独", "一次性交互")
	context    = datautil.NewTargetValues("context", "上下文", "上下文交互", "语境", "语境交互")
	typeActive = datautil.NewTargetValues("type", "信息类型", "类型")
)

//帮助功能，用于显示全部命令
type HelpInteract struct {
}

func NewHelpInteract() interact.SingleMessageInteract {
	return HelpInteract{}
}

func (interact HelpInteract) GetCommandName() []string {
	return []string{"help", "帮助", "功能"}
}
func (interact HelpInteract) RespondSource() []constdata.MessageType {
	return []constdata.MessageType{
		constdata.GroupMessage,
	}
}

func (i HelpInteract) EnterMessage(
	extraCmd datautil.MutliToOneMap,
	data structs.Message,
	repChan chan messagetargets.MessageTarget) {
	var msg, _ = sourceHandle.GetGoupSoucreMessage(data.Source)

	if len(extraCmd) == 0 {
		var d []request.H

		d = append(d, request.H{"type": string(constdata.Plain),
			"text": "收到信息：帮助\n"})

		d = append(d, request.H{
			"type": string(constdata.Plain),
			"text": "单次交互命令:\n",
		})

		for _, v := range interact.GetSingleCommand() {
			d = append(d, request.H{
				"type": string(constdata.Plain),
				"text": " #" + v + "\n",
			})
		}

		d = append(d, request.H{
			"type": string(constdata.Plain),
			"text": "上下文交互命令\n",
		})

		for _, v := range interact.GetContextCommand() {
			d = append(d, request.H{
				"type": string(constdata.Plain),
				"text": " # " + v + "\n",
			})
		}

		d = append(d, request.H{
			"type": string(constdata.Plain),
			"text": "信息类型单次响应\n",
		})

		for _, v := range interact.GetChainSingleCommand() {
			d = append(d, request.H{
				"type": string(constdata.Plain),
				"text": "类型： " + v + "\n",
			})
		}

		var da = messagetargets.NewGroupTarget(msg.GroupId, d)

		repChan <- da
	} else {
		cmdName, nok := extraCmd.Get("name", "指令", "命令")
		cmdName = strings.ToLower(cmdName)
		sourceName, _ := extraCmd.GetWithDefault("single", "type", "来源", "类型")

		if nok {
			if single.Match(sourceName) {
				c, ok := interact.GetSingleInteract(cmdName)
				if ok {
					repChan <- messagetargets.NewSingleTextGroupTarget(msg.GroupId, c().GetUseage())
					return
				}
			} else if context.Match(sourceName) {
				c, ok := interact.GetContextInteract(cmdName)
				if ok {
					repChan <- messagetargets.NewSingleTextGroupTarget(msg.GroupId, c().GetUseage())
					return
				}
			}else if typeActive.Match(sourceName){
				c, ok := interact.GetSingleInteract(cmdName)
				if ok {
					repChan <- messagetargets.NewSingleTextGroupTarget(msg.GroupId, c().GetUseage())
					return
				}
			}
			repChan <- messagetargets.NewSingleTextGroupTarget(msg.GroupId, "指令未找到！")
			return
		}
		repChan <- messagetargets.NewSingleTextGroupTarget(msg.GroupId, "指令不完整！name 和 type 要同时指定")
		return
	}
}

func (HelpInteract) GetUseage() string {
	return "#help|#帮助 : 返回当前机器人的全部指令\n" +
		"额外指令：\n" +
		" name|指令|命令=[指令名称] \n  -> 标记选择的指令名称\n" +
		" type|来源|类型=[single|context] \n  -> 选择指令来源【单次交互|上下文交互】"
}
