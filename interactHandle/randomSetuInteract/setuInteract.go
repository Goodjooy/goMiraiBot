package randomsetuinteract

import (
	"goMiraiQQBot/lib/client"
	"goMiraiQQBot/lib/constdata"
	datautil "goMiraiQQBot/lib/dataUtil"
	interactprefab "goMiraiQQBot/lib/interactPrefab"
	"goMiraiQQBot/lib/messageHandle/interact"
	messagetargets "goMiraiQQBot/lib/messageHandle/messageTargets"
	"goMiraiQQBot/lib/messageHandle/structs"
	"strconv"
)

var apiKey string
var targetURL string

var (
	setuCount     = datautil.NewTargetValues("count", "数量")
	setuR18       = datautil.NewTargetValues("r18", "R18")
	setuThumbnail = datautil.NewTargetValues("thumbnail", "缩略图")
)

type SetuInteract struct {
	interactprefab.InteractPerfab
}

func NewSetuInteract() interact.FullSingleInteract {
	perfab := interactprefab.
		NewInteractPerfab().
		AddActivateSigns("setu").
		AddActivateSigns("随机涩图").
		AddActivateSigns("涩图").
		AddActivateSource(constdata.GroupMessage).
		AddActivateSource(constdata.FriendMessage).
		AddInitFn(func() {
			setuCfg, ok := client.GetExtraConfig("setu")
			if !ok {
				apiKey = ""
				targetURL = "https://api.lolicon.app/setu/"
			} else {
				setu := setuCfg.(map[interface{}]interface{})
				if key, ok := setu["apiKey"]; ok {
					apiKey = key.(string)
				} else {
					apiKey = ""
				}
				if key, ok := setu["url"]; ok {
					targetURL = key.(string)
				} else {
					targetURL = "https://api.lolicon.app/setu/"
				}
			}
		}).
		SetUseage(`随机涩图功能，每日额度300`).
		BuildPtr()

	return &SetuInteract{InteractPerfab: *perfab}
}

func (setu *SetuInteract) EnterMessage(
	extraCmd datautil.MutliToOneMap,
	data structs.Message,
	repChan chan messagetargets.MessageTarget) {

	extraCmd.SetNoNameCmdOrder(setuCount, setuThumbnail, setuR18)
	setuS, _ := extraCmd.GetWithDefault("1", setuCount...)
	setuC, err := strconv.Atoi(setuS)
	if err != nil {
		setuC = 1
	}
	thumbnailS, _ := extraCmd.GetWithDefault("true", setuThumbnail...)
	thumbnailC, err := strconv.ParseBool(thumbnailS)
	if err != nil {
		thumbnailC = true
	}

	randomSetu(
		false,
		uint(setuC),
		thumbnailC,
		data.Source,
		repChan,
	)
}
