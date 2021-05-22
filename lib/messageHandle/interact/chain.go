package interact

import (
	"goMiraiQQBot/lib/messageHandle/structs"
	"strings"
)

func ChainBaseInteractController() InteractController {
	return NewInteractContorller(
		0,
		func(m structs.Message) (Command, bool) {
			return chainStructGet(m.ChainInfoList)
		},
	)
}

func chainStructGet(chains []structs.MessageChainInfo) (Command, bool) {
	var cmd Command

	if len(chains) >= 1 {
		cmd.mainCmd = strings.ToLower(string(chains[0].MessageType))
		return cmd, true
	}
	return cmd, false
}
