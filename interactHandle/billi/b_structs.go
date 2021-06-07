package billi

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

type BiliRes struct {
	Code int32  `json:"code"`
	Mes  string `json:"message"`
	Ttl  int32  `json:"ttl"`

	Data interface{} `json:"data"`
}

type UserInfoRes struct {
	Code int32  `json:"code"`
	Mes  string `json:"message"`
	Ttl  int32  `json:"ttl"`

	Data BilUInfo `json:"data"`
}

type BilUInfo struct {
	Face string `json:"face"`
	Name string `json:"name"`
}

func GetUPerInfo(uid string) (UserInfoRes, error) {
	values := url.Values{}
	values.Add("mid", uid)
	values.Add("jsonp", "jsonp")

	targetURL := url.URL{
		Scheme:   "https",
		Host:     "api.bilibili.com",
		Path: "/x/space/acc/info",
		RawQuery: values.Encode(),
	}

	res, err := http.Get(targetURL.String())
	if err != nil {
		return UserInfoRes{}, err
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return UserInfoRes{}, err
	}

	var user UserInfoRes
	err = json.Unmarshal(data, &user)
	if err != nil {
		return user, err
	}

	return user, nil
}
