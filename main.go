package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"goMiraiQQBot/constdata"
	"goMiraiQQBot/messageHandle/interact"
	"goMiraiQQBot/messageHandle/interactHandle"
	hentaiimageinteract "goMiraiQQBot/messageHandle/interactHandle/hentaiImageInteract"
	messagetargets "goMiraiQQBot/messageHandle/messageTargets"
	"goMiraiQQBot/messageHandle/sourceHandle"
	s "goMiraiQQBot/messageHandle/structs"
	"goMiraiQQBot/request"
	"goMiraiQQBot/request/structs"
	"goMiraiQQBot/request/structs/message"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

const authKey = "INITKEYy2euAf0E"

var addr = flag.String("addr", "0.0.0.0:8080", "http service address")

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
	if err != nil || resInterface.Code != constdata.Normal {
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
	if err != nil || res.Code != constdata.Normal {
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

	resMessage := make(chan messagetargets.MessageTarget)
	reqMessage := make(chan s.Message)

	sourceHandle.InitMessageSourceHandle()

	interact.AddSingleInteract(interactHandle.NewHelpInteract)
	interact.AddSingleInteract(interactHandle.NewAboutInteract)

	interact.AddContextInteract(hentaiimageinteract.NewHentaiImageSearchInteract)

	
	interact.InitInteractHandle(reqMessage, resMessage)

	done := make(chan struct{})
	//load message

	messageReciver := func() {
		defer close(done)

		for {
			var f message.MessageMapRespond
			_,info,err:=socket.ReadMessage()
			if err != nil {
				log.Fatal("readMessageFatal: ", err)
				return
			}
			err=json.Unmarshal(info,&f.Data)
			if err != nil {
				log.Fatal("Unmarshal Json Failure",err)
				continue
			}


			msg, err := s.FromMessageRespondData(f)
			if err != nil {
				log.Fatal(err)
				continue
			}

			reqMessage <- msg
			log.Print("send Message To Handle")
		}
	}

	messageSender := func() {
		for {
			select {
			case data, ok := (<-resMessage):
				if ok {

					var result message.MessageSendRespond
					err := request.PostWithTargetRespond(string(data.GetTargetPort()), data.GetSendContain(resInterface.Session),&result)
					if err != nil {
						log.Fatalf("Send Message Fail %v", err)
						continue
					}
					if result.Code!=constdata.Normal{
						log.Fatal("Bad Respond Code: ",(result.Code))
					}else {
						log.Printf("Success Send Message! messageId:%v",result.MessageId)
					}
				}
			}
		}
	}

	go messageReciver()
	go messageReciver()

	go messageSender()
	go messageSender()
	go messageSender()
	go messageSender()

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
				socket.WriteMessage(websocket.TextMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				return
			}
		}
	}
}
