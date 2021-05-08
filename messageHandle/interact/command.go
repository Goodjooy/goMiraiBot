package interact

import (
	"goMiraiQQBot/constdata"
	datautil "goMiraiQQBot/dataUtil"

	"goMiraiQQBot/messageHandle/structs"
	"regexp"
	"strings"
)

var cmdPattern = regexp.MustCompile(`^(#\s*(\S+))\s*`)
var extraCmdPattern = regexp.MustCompile(`(?:\s)(\S+)[=:：](\S+)`)

type command struct {
	mainCmd  string
	extraCmd datautil.MutliToOneMap
}

func commandGet(msgChain []structs.MessageChainInfo) (command, bool) {
	cmd := command{}
	msg := commandLoad(msgChain)
	//以#开头,有命令
	if strings.HasPrefix(msg, "#") {
		matchResult := cmdPattern.FindStringSubmatch(msg)
		cmd.mainCmd = strings.ToLower(matchResult[2])

		extraCmd := msg[len(matchResult[1]):]

		extraMatch := extraCmdPattern.FindAllStringSubmatch(extraCmd, -1)
		cmd.extraCmd = datautil.NewMutliToOneMap()
		for _, v := range extraMatch {
			cmd.extraCmd[v[1]] = v[2]
		}

		return cmd, true
	}
	return cmd, false
}

func commandLoad(msgCHain []structs.MessageChainInfo) string {
	var cmd string
	for _, v := range msgCHain {
		if v.MessageType == constdata.Plain {
			cmd += v.Data["text"].(string)
		}
	}
	return cmd
}

func atCommandLoad(msgChain []structs.MessageChainInfo)bool {
	//TODO: 监控是否@自己
	return false
}