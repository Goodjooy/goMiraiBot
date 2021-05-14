package client

import (
	"encoding/json"
	"fmt"
	"goMiraiQQBot/constdata"
	messagetargets "goMiraiQQBot/messageHandle/messageTargets"
	"goMiraiQQBot/messageHandle/structs"
	"goMiraiQQBot/request"
	"goMiraiQQBot/request/structs/message"
	"io/ioutil"
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

type WSHolder struct {
	Conn *websocket.Conn
}

func MessageReader(msgSocket *websocket.Conn, reqMsg chan structs.Message, session Session, conn *WSHolder, url url.URL) bool {
	var f message.MessageMapRespond
	_, msgReader, err := msgSocket.NextReader()
	if err != nil {
		log.Print("Read Message | Get Message Failure: ", err)
		//cnn, _ := TryReDialWebSocket(EstablishMessageHandleWebSocket, 6, session, url)
		//conn.Conn = cnn
		return false
	}
	var data []byte
	data, err = ioutil.ReadAll(msgReader)
	if err != nil {
		log.Print("Read Message | Read Data From Reader Fauiure: ", err)
		return true
	}

	err = json.Unmarshal(data, &f.Data)
	if err != nil {
		log.Print("Read Message | Unmarshal Json Failure", err)
		return true
	}

	msg, err := structs.FromMessageRespondData(f)
	if err != nil {
		log.Print("Read Message | Struct Transfrom Error: ", err)
		return true
	}

	log.Printf("Read Message | Accept Message Success! Source: \n%+v\n`%+v`",
		msg.Source,
		chainRander(msg.ChainInfoList))
	reqMsg <- msg

	return true
}

func MessageSender(data messagetargets.MessageTarget, session Session) {
	if data == nil {
		return
	}

	log.Printf("Handle Message Send: %v", data.GetSendMessage())
	var result message.MessageSendRespond
	err := request.PostWithTargetRespond(string(data.GetTargetPort()), data.GetSendContain(string(session)), &result)
	if err != nil {
		log.Printf("Send Message Failure: %v", err)
		return
	}
	if result.Code != constdata.Normal {
		log.Print("Bad Respond Code: ", (result.Code))
	} else {
		log.Printf("Success Send Message! messageId: %v", result.MessageId)
	}
}

func MessageReaderHolder(done chan struct{},
	reqMsg chan structs.Message,
	session Session,
	msgSocket *websocket.Conn,
	url url.URL,
	conn *WSHolder) {
	//退出函数关闭阻塞挂起

	for {
		ok := MessageReader(msgSocket, reqMsg, session, conn, url)
		if !ok {
			return
		}
	}
}

func MessageSenderHolder(done chan struct{}, targetChan chan messagetargets.MessageTarget, session Session) {
	for {
		select {
		case <-done:
			log.Print("Cilent Stop, Close Message Sender")
			return
		case data, ok := (<-targetChan):
			if ok {
				MessageSender(data, session)
			}
		}
	}
}

func chainRander(chains []structs.MessageChainInfo) string {
	var info string

	for _, v := range chains {
		if v.MessageType == constdata.At {
			info += fmt.Sprintf("[AT-%v-%v] ", uint64(v.Data["target"].(float64)), v.Data["display"].(string))
		} else if v.MessageType == constdata.Plain {
			info += v.Data["text"].(string)
		} else {
			info += fmt.Sprintf("[TYPE:%v]", v.MessageType)
		}
	}
	return info
}
