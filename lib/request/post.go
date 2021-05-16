package request

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

type H map[string]interface{}

const contextType = "application/json"
const urlBody = "http://0.0.0.0:8080"

func Post(path string, body interface{}) (interface{}, error) {
	var f interface{}
	err:=PostWithTargetRespond(path,body,&f)
	if err != nil {
		return nil, err
	}
	return f,nil
}

func PostWithTargetRespond(path string, body interface{},resStructPtr interface{})error{
	 //data transfrom
	//to json
	jsonByte, err := json.Marshal(body)
	if err != nil {
		return  errors.New("Marshal Failure : " + err.Error())
	}
	//byte Reader
	jsonReder := bytes.NewReader(jsonByte)

	//sendMessage
	res, err := http.Post(urlBody+path, contextType, jsonReder)
	if err != nil {
		return  errors.New("Bad Respond : " + err.Error())
	}
	//recive message
	respondByte, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return  errors.New("Failure Read RespondBody : " + err.Error())
	}
	//transfrom respond body
	
	err = json.Unmarshal(respondByte, &resStructPtr)
	if err != nil {
		return  errors.New("Transfrom Respond Body To Json Fail : " + err.Error()+"\n"+ string(respondByte))
	}

	if res.StatusCode != http.StatusOK {
		return  errors.New("Bad Respond Status : " + res.Status)
	}


	
	return  nil
}