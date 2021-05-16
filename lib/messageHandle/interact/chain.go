package interact

import (
	"goMiraiQQBot/lib/messageHandle/structs"
	"strings"
)


func chainStructGet(chains []structs.MessageChainInfo) (Command,bool) {
	var cmd Command

	if len(chains) >= 1 {
		cmd.mainCmd = strings.ToLower(string(chains[0].MessageType))
		return cmd,true
	}
	return cmd ,false
}
