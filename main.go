package main

import (
	"flag"
	"fmt"
	"goMiraiQQBot/request"
	"goMiraiQQBot/request/structs"
	"goMiraiQQBot/request/structs/message"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

const authKey = "INITKEYy2euAf0E"

var addr = flag.String("addr", "0.0.0.0:8080", "http service address")

var messagechannal=make(chan message.MessageData, 128)

func main() {
	flag.Parse()
	log.SetFlags(1)

	//close workser
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	//get session
	var resInterface structs.AuthRespond
	err := request.PostWithTargetRespond(
		"/auth",
		request.H{
			"authKey": authKey,
		},
		&resInterface)
	if err != nil || resInterface.Code != request.Normal {
		log.Fatal(err)
		return
	}

	log.Println("send auth key")

	//verify qq
	verifyRequestBody := structs.VerifyQQRequest{
		QQ:         3628862306,
		SessionKey: resInterface.Session,
	}
	var res structs.VerifyRespond

	err = request.PostWithTargetRespond("/verify", verifyRequestBody, &res)
	if err != nil || res.Code != request.Normal {
		log.Fatal(err)
		return
	}
	log.Println("Verify QQ")


	//ws url
	u := url.URL{Scheme: "ws", Host: *addr, Path: "/message", RawQuery: fmt.Sprintf("sessionKey=%s", resInterface.Session)}
	socket, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
		return
	}
	log.Println("Bot Started")

	defer socket.Close()

	done := make(chan struct{})
	//load message

	go func() {
		defer close(done)

		for {
			var f message.MessageData
			err := socket.ReadJSON(&f)
			if err != nil {
				log.Fatal("readMessageFatal: ", err)
				return
			}

			if f.Type == request.GroupMessage {
				go func() {
					chain := f.MessageChain
					messageChain := chain[1]
					var code string = ""
					if messageChain["type"].(string) == string(request.Plain) {
						text := messageChain["text"].(string)

						if !strings.HasPrefix(text, "#") {
							return
						}
						code = strings.Replace(text, "#", "", 1)
					} else if messageChain["type"].(string) == string(request.Image) {

						code = fmt.Sprintf("收到图片")
					} else {
						return
					}

					msg := message.GroupMessageRequest{
						Session: resInterface.Session,
						Target:  f.Sender.GroupIn.Id,
						Clain: []request.H{
							{
								"type": string(request.Plain),
								"text": fmt.Sprintf("收到信息！\n 来自群[%v(%v)]\n发送者[%v(%v)]\n`%v`",
									f.Sender.GroupIn.Name, f.Sender.GroupIn.Id,
									f.Sender.MemberName, f.Sender.Id, code),
							}},
					}
					if messageChain["type"].(string) == string(request.Image) {
						msg.Clain = append(msg.Clain, request.H{
							"type": string(request.Image),
							"url":  messageChain["url"].(string),
						})
					}
					var res message.MessageSendRespond
					err = request.PostWithTargetRespond("/sendGroupMessage", msg, &res)
					if err != nil {
						log.Fatal(err)
						return
					} else {
						log.Printf("send message success :group[%v(%v)]", f.Sender.GroupIn, f.Sender.GroupIn.Id)
					}
				}()

			}

		}

	}()
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case <-interrupt:
			{
				log.Print("interrupt")
				//release sesstion
				func() {
					releaseSessionBody := structs.VerifyQQRequest{
						QQ:         3628862306,
						SessionKey: resInterface.Session,
					}
					var res structs.VerifyRespond
			
					err := request.PostWithTargetRespond("/release", releaseSessionBody, &res)
					if err != nil {
						log.Fatal(err)
					}
					log.Print("exit release session")
				}()
				socket.WriteMessage(websocket.TextMessage,websocket.FormatCloseMessage(websocket.CloseNormalClosure,""))				
				return
			}
		}
	}
}
