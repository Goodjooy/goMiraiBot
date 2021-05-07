package hentaiimageinteract

import (
	"fmt"
	"goMiraiQQBot/constdata"
	messagetargets "goMiraiQQBot/messageHandle/messageTargets"
	"goMiraiQQBot/request"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"sync"
	"time"
)

const api = "https://ascii2d.net"
const uri = "https://ascii2d.net/search/uri"

var crsfParamPattern = regexp.MustCompile(`<meta name="csrf-param" content="([a-zA-Z_]+)" />`)
var crsfTokenPattern = regexp.MustCompile(`<meta name="csrf-token" content="([^<>]+)" />`)

var lastTimeMap sync.Map

type crsfToken struct {
	prarm string
	token string
}

func searchHandle(groupId uint64, imageURL string, resChan chan messagetargets.MessageTarget) {
	resChan <- messagetargets.NewSingleTextGroupTarget(groupId, "正在以图搜图...")

	//时间限制，30s内只进行一次
	lastTime, ok := lastTimeMap.Load(groupId)
	if !ok {
		lastTimeMap.Store(groupId, time.Now())
	} else {
		last := lastTime.(time.Time)
		if last.Add(30 * time.Second).After(time.Now()) {
			resChan <- messagetargets.NewSingleTextGroupTarget(groupId, "请求搜图间隔过短[<30s]")
			return
		} else {
			lastTimeMap.Store(groupId, time.Now())
		}
	}

	crsf, err := getApiCrsfToken()
	if err != nil {
		log.Fatal(err)
		resChan <- messagetargets.NewSingleTextGroupTarget(groupId, fmt.Sprintf("获取搜图token失败: %v", err))
		return
	}
	results, err := postForm(crsf, imageURL)
	if err != nil {
		log.Fatal(err)
		resChan <- messagetargets.NewSingleTextGroupTarget(groupId, fmt.Sprintf("发起图片搜索失败：%v", err))
		return
	}

	resChan <- messagetargets.NewSingleTextGroupTarget(groupId, fmt.Sprintf("搜索到%v个结果,将发送前4个结果\n来源：%v", len(results), api))
	for _, v := range results[:4] {
		resChan <- messagetargets.NewGroupTarget(groupId, v.toChains())
	}

}

func getApiCrsfToken() (crsfToken, error) {
	res, err := http.Get(api)
	if err != nil {
		return crsfToken{}, err
	}

	htmlHead, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return crsfToken{}, err
	}

	paramName := crsfParamPattern.FindStringSubmatch(string(htmlHead))[1]
	token := crsfTokenPattern.FindStringSubmatch(string(htmlHead))[1]

	return crsfToken{
		prarm: paramName,
		token: token,
	}, nil

}

func postForm(crsf crsfToken, imageURL string) ([]searchResult, error) {
	values := url.Values{}

	values.Set("utf8", "✓")
	values.Set(crsf.prarm, crsf.token)
	values.Set("uri", imageURL)

	res, err := http.PostForm(uri, values)
	if err != nil {
		return nil, err
	}
	bodyText, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	thumbnails := thumbnailPicPattern.FindAllStringSubmatch(string(bodyText), -1)
	linkDatas := metaDataPatern.FindAllStringSubmatch(string(bodyText), -1)

	size := len(linkDatas)
	var results []searchResult

	for i := 0; i < size; i++ {
		t := thumbnails[i]
		l := linkDatas[i]
		s := searchResult{
			thumbnailPic: t[1],
			workURL:      l[1],
			workName:     l[2],
			userURL:      l[3],
			userName:     l[4],
			source:       l[5],
		}
		results = append(results, s)
	}

	return results, nil
}

type searchResult struct {
	thumbnailPic string

	workURL  string
	workName string

	userURL  string
	userName string

	source string
}

var metaDataPatern = regexp.MustCompile(
	`\s*<h6>\s*` +
		`\s*<img class="to-link-icon" src="(?:[^<>]+)" alt="[^<>]+" width="\d+" height="\d+" />\s*` +
		`\s*<a target="_blank" rel="noopener" href="([^<>]+)">([^<>]+)</a>\s*` +
		`\s*<a target="_blank" rel="noopener" href="([^<>]+)">([^<>]+)</a>\s*` +
		`\s*<small>\s*` +
		`\s*([^<>]+)\s*` +
		`\s*</small>\s*` +
		`\s*</h6>\s*`)
var thumbnailPicPattern = regexp.MustCompile(`<img loading="lazy" src="([^<>]+)" alt="[^<>]+" width="\d+" height="\d+" />`)

func (r searchResult) toChains() []request.H {
	return []request.H{
		{
			"type": constdata.Plain.String(),
			"text": "搜索到结果\n",
		},
		{
			"type": constdata.Image.String(),
			"url":  api + r.thumbnailPic,
		},
		{
			"type": constdata.Plain.String(),
			"text": fmt.Sprintf("作品url:`%v`\n作品名称：`%v`\n", r.workURL, r.workName),
		}, {
			"type": constdata.Plain.String(),
			"text": fmt.Sprintf("作者主页:`%v`\n作者名称：`%v`\n", r.userURL, r.userName),
		}, {
			"type": constdata.Plain.String(),
			"text": "---------------\n",
		},
	}
}
