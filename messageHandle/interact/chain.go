package interact

import (
	"goMiraiQQBot/messageHandle/structs"
	"strings"
)


func chainStructGet(chains []structs.MessageChainInfo) (command,bool) {
	var cmd command

	if len(chains) >= 1 {
		cmd.mainCmd = strings.ToLower(string(chains[0].MessageType))
		return cmd,true
	}
	return cmd ,false
}
