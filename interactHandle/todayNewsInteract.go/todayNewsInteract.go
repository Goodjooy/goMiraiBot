package todaynewsinteractgo

import (
	"encoding/json"
	"goMiraiQQBot/lib/constdata"
	datautil "goMiraiQQBot/lib/dataUtil"
	interactprefab "goMiraiQQBot/lib/interactPrefab"
	"goMiraiQQBot/lib/messageHandle/interact"
	messagetargets "goMiraiQQBot/lib/messageHandle/messageTargets"
	"goMiraiQQBot/lib/messageHandle/structs"
	"io/ioutil"
	"log"
	"net/http"
)
type dailyMessage struct {
	Code uint `json:"code"`
	Message string `json:"msg"`

	ImageURL string `json:"imageUrl"`
}


var targetURL string ="http://dwz.2xb.cn/zaob"
type TodayNewsInteract struct {
	interactprefab.InteractPerfab
}

func NewTodyNewsInteract() interact.FullSingleInteract {
	return &TodayNewsInteract{
		InteractPerfab: interactprefab.NewInteractPerfab().
			AddActivateSigns("新闻", "news").
			AddActivateSource(constdata.GroupMessage).
			SetUseage("60秒看世界").
			AddInitFn(func() {

			}).
			Build(),
	}
}

func (tn *TodayNewsInteract) EnterMessage(
	extraCmd datautil.MutliToOneMap,
	data structs.Message,
	repChan chan messagetargets.MessageTarget) {

	res,err:=http.Get(targetURL)
	//request error
	if err != nil {
		log.Printf("Get Daliy Image Url Failure: %v",err)
		return
	}
	//read request body error
	resData ,err:=ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("Read Respond Body Error: %v",err)
		return
	}

	//unmarsal to Json
	var dailyInfo dailyMessage
	err=json.Unmarshal(resData,&dailyInfo)
	if err != nil {
		log.Printf("Unmarshal to Json Error : %v",err)
		return
	}

	repChan<-messagetargets.SourceTarget(
		data.Source,
		structs.NewImageChain(dailyInfo.ImageURL),
	)
}
