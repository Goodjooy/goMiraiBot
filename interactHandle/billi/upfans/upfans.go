package upfans

import (
	"encoding/json"
	"fmt"
	"goMiraiQQBot/interactHandle/billi"
	"goMiraiQQBot/lib/constdata"
	datautil "goMiraiQQBot/lib/dataUtil"
	interactprefab "goMiraiQQBot/lib/interactPrefab"
	"goMiraiQQBot/lib/messageHandle/interact"
	messagetargets "goMiraiQQBot/lib/messageHandle/messageTargets"
	"goMiraiQQBot/lib/messageHandle/sourceHandle"
	"goMiraiQQBot/lib/messageHandle/structs"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
)

var lastTimeMap sync.Map
var user_id = datautil.NewTargetValues("uid", "用户id", "id")




type Upfans struct {
	interactprefab.InteractPerfab
}

func NewUpFansInteract() interact.FullSingleInteract {
	return &Upfans{
		InteractPerfab: interactprefab.
			NewInteractPerfab().
			AddActivateSigns("up粉丝", "up-fans").
			AddActivateSource(constdata.GroupMessage).
			SetUseage("# up粉丝 uid:<用户uid>\n获取目标用户的粉丝数量").
			AddInitFn(func() {}).
			Build(),
	}
}

// EnterMessage 响应信息
func (*Upfans) EnterMessage(
	extraCmd datautil.MutliToOneMap,
	data structs.Message,
	repChan chan messagetargets.MessageTarget) {
	extraCmd.SetNoNameCmdOrder(user_id)

	uid, ok := extraCmd.Get(user_id.GetAll()...)
	if !ok {
		repChan <- messagetargets.SourceTarget(data.Source, structs.
			NewTextChain("请提供UP的ID\n如果要查询的UP的ID为114145，指令如下\n #up粉丝 uid:114145"))
		return
	} else {
		go get_up_fans(uid, data.Source, repChan)
	}

}

func get_up_fans(uid string,
	source sourceHandle.MessageSource,
	repChan chan messagetargets.MessageTarget) {
	var last_time time.Time
	last_time_i, ok := lastTimeMap.Load(uid)
	if ok {
		last_time = last_time_i.(time.Time)
		now := time.Now()
		if last_time.Add(60 * time.Second).After(now) {
			delta := last_time.Add(60 * time.Second).Sub(now)
			repChan <- messagetargets.SourceTarget(source,
				structs.NewTextChain(fmt.Sprintf("请求间隔过短，请求间隔应该大于60s(%vs)", int(delta/time.Second))))
			return
		} else {
			lastTimeMap.Store(uid, now)
		}
	}else{
		lastTimeMap.Store(uid, time.Now())
	}

	values := url.Values{}
	values.Add("vmid", uid)
	values.Add("jsonp", "jsonp")
	targetURL := url.URL{
		Scheme: "https",
		Host:     "api.bilibili.com",
		Path:     "/x/relation/stat",
		RawQuery: values.Encode(),
	}
	user,err:=billi.GetUPerInfo(uid)
	if err != nil {
		user.Data.Name=uid	
	}

	res, err := http.Get(targetURL.String())
	if err != nil {
		log.Printf("Get UPer Fans Error: %v", err)
		return
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("Read Uper Fans Respond Body Error: %v", err)
		return
	}

	var bRes billi.BiliRes
	err = json.Unmarshal(data, &bRes)
	if err != nil {
		log.Printf("Unmarshal Uper Fans Respond Body Error: %v", err)
		return
	}

	fansData, ok := bRes.Data.(map[string]interface{})["follower"]
	if !ok {
		log.Printf("No Uper Fans Mesaage found in respond body")
		return
	}
	fans := fmt.Sprintf("%.0f", fansData.(float64))

	repChan <- messagetargets.SourceTarget(source,
		structs.NewAtChain(source.GetSenderID(), ""),
		structs.NewImageChain(user.Data.Face),
		structs.NewTextChain(fmt.Sprintf("\n 用户为 `%v` 的粉丝数量为：\n%v", user.Data.Name, fans)))

}
