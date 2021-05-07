package hentaiimageinteract

import (
	"goMiraiQQBot/constdata"
	"goMiraiQQBot/messageHandle/structs"
)

func foundTargeImage(chains []structs.MessageChainInfo) (string, bool) {
	for _, v := range chains {
		//第一个图片
		if v.MessageType == constdata.Image {
			data := v.Data

			if URL, ok := data["url"]; ok {
				return URL.(string), true
			}
			break
		}
	}
	return "", false
}

func findCancelSign(chains []structs.MessageChainInfo) bool {
	for _, v := range chains {
		//第一个图片
		if v.MessageType == constdata.Plain {
			data := v.Data

			if text, ok := data["text"]; ok {
				if text.(string) == "取消" {
					return true
				}
			}
			break
		}
	}
	return false
}
