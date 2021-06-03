package interact

import (
	"goMiraiQQBot/lib/constdata"
	datautil "goMiraiQQBot/lib/dataUtil"

	"goMiraiQQBot/lib/messageHandle/structs"
	"regexp"
	"strings"
)

var cmdPattern = regexp.MustCompile(`^(#\s*(\S+))\s*`)
var extraCmdPattern = regexp.MustCompile(`(?:\s)(\S+)[=:：](\S+)`)
var extraNoNameCmdPattern = regexp.MustCompile(`(?:\s)(\S+)`)

type Command struct {
	mainCmd  string
	extraCmd datautil.MutliToOneMap
}

func (c *Command) GetExtraCmd() datautil.MutliToOneMap {
	return c.extraCmd
}

func NewCommandBaseInteractGroup() InteractController {
	return NewInteractContorller(
		1,
		func(m structs.Message) (Command, bool) {
			return CommandGet(m.ChainInfoList, cfg.GetQQId())
		})
}

func CommandGet(msgChain []structs.MessageChainInfo, botQQ uint64) (Command, bool) {
	cmd := Command{}
	msg := commandLoad(msgChain, botQQ)
	//以#开头,有命令
	if strings.HasPrefix(msg, "#") {
		matchResult := cmdPattern.FindStringSubmatch(msg)
		if matchResult == nil || len(matchResult) < 3 {
			return cmd, false
		}
		cmd.mainCmd = strings.ToLower(matchResult[2])

		extraCmd := msg[len(matchResult[1]):]

		extraMatch := extraCmdPattern.FindAllStringSubmatch(extraCmd, -1)
		cmd.extraCmd = datautil.NewMutliToOneMap()
		for _, v := range extraMatch {
			cmd.extraCmd.Put(v[1], v[2])
		}

		extraNoNameCmd := extraNoNameCmdPattern.FindAllStringSubmatch(extraCmd, -1)
		for _, v := range extraNoNameCmd {
			cmd.extraCmd.PutNoName(v[1])
		}

		return cmd, true
	}
	return cmd, false
}

func commandLoad(msgCHain []structs.MessageChainInfo, botQQ uint64) string {
	var cmd string
	if msgCHain[0].MessageType == constdata.At {
		qqId := uint64(msgCHain[0].Data["target"].(float64))
		if qqId == uint64(botQQ) {
			cmd += "#"
			msgCHain = msgCHain[1:]
		}
	}

	for _, v := range msgCHain {
		if v.MessageType == constdata.Plain {
			cmd += v.Data["text"].(string)
		}
	}
	return cmd
}
