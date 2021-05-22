package randomsetuinteract

import (
	"encoding/json"
	"fmt"
	messagetargets "goMiraiQQBot/lib/messageHandle/messageTargets"
	"goMiraiQQBot/lib/messageHandle/sourceHandle"
	"goMiraiQQBot/lib/messageHandle/structs"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

func randomSetu(
	r18 uint,
	num uint,
	thumbnail bool,
	source sourceHandle.MessageSource,
	msgRepChan chan messagetargets.MessageTarget,
) {
	values := url.Values{}

	values.Add("apikey", apiKey)
	values.Add("r18", strconv.FormatUint(uint64(r18), 10))
	values.Add("num", strconv.FormatUint(uint64(num), 10))
	values.Add("size1200", strconv.FormatBool(thumbnail))

	res, err := http.Get(fmt.Sprintf("%s?%s", targetURL, values.Encode()))
	//request error
	if err != nil {
		log.Printf("conncet to Setu service Failure: %v", err)
		return
	}
	//data read error
	var data []byte
	data, err = ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("read respond Body Failure: %v", err)
		return
	}

	//unmarchal To Json error
	var setuData setuRes
	err = json.Unmarshal(data, &setuData)
	if err != nil {
		log.Printf("Unmarshal to Json Failure: %v", err)
		return
	}

	//code config
	if setuData.Code != 0 {
		log.Printf("Bad Respond Code: %v | %v", setuData.Code, setuData.Msg)
		msgRepChan <- messagetargets.SourceTarget(source,
			structs.NewTextChain(
				fmt.Sprintf(
					"获取涩图失败：\n  错误码：%v\n  错误信息：%v",
					setuData.Code,
					setuData.Msg,
				),
			),
		)
		return
	}

	//success
	msgRepChan <- messagetargets.SourceTarget(source,
		structs.NewTextChain(fmt.Sprintf(
			"获取涩图成功，即将发送%v张", setuData.Count)))
	msgRepChan <- messagetargets.SourceTarget(source,
		structs.NewTextChain(
			fmt.Sprintf("本日剩余接口调用次数：%v", setuData.Quota),
		))

	for _, v := range setuData.Setus {
		setuInfo := fmt.Sprintf("作品ID:%v\n图片序列：%v\n作者：%v",
			v.Pid, v.PicNum, v.Author)

		msgRepChan <- messagetargets.SourceTarget(source,
			structs.NewTextChain(setuInfo),
			structs.NewImageChain(v.URL))
	}

}
