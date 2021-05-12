package client

import (
	"fmt"
	"goMiraiQQBot/messageHandle/interact"
	messagetargets "goMiraiQQBot/messageHandle/messageTargets"
	"goMiraiQQBot/messageHandle/sourceHandle"
	mStruct "goMiraiQQBot/messageHandle/structs"
	"goMiraiQQBot/request"
	"goMiraiQQBot/request/structs"
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

type BotWsClient struct {
	rootURL url.URL

	msgSocket   *websocket.Conn
	eventSocket *websocket.Conn

	session Session

	Done chan struct{}

	msgGetChan chan mStruct.Message
	msgReq     chan messagetargets.MessageTarget

	config Config
}

func NewQQBotClient(configPath string) BotWsClient {
	client := BotWsClient{
		Done:       make(chan struct{}),
		msgGetChan: make(chan mStruct.Message),
		msgReq:     make(chan messagetargets.MessageTarget),
	}

	config := LoadConfig(configPath)
	client.config = config

	return client
}

func (client *BotWsClient) Run() error {
	log.Print("Starting QQ Bot")
	config := client.config

	log.Print("Init Souce Handle...")
	sourceHandle.InitMessageSourceHandle()
	log.Print("Souce Handle Init Success!")

	log.Print("Init Interact Handle")
	interact.InitInteractHandle(client.msgGetChan, client.msgReq)
	log.Print("Interact Handle Init Success")

	session, err := AuthQQKey(config)
	if err != nil {
		log.Fatal("Init Bot Clent Faulure | Send Auth Key: ", err)
		return err
	} else {
		client.session = session
	}
	err = VerifyQQ(session, config)
	if err != nil {
		log.Fatal("Init Bot Clent Faulure | Verify QQ: ", err)
		return err
	}

	log.Print("Establish Websocket Connect To Bot")
	client.rootURL = url.URL{Host: fmt.Sprintf("%v:%v", config.Server.Host, config.Server.Port)}

	client.msgSocket, err = EstablishMessageHandleWebSocket(client.rootURL, session)
	if err != nil {
		log.Fatal("Init Bot Clent Faulure | Establish Websocket at `/message` : ", err)
		return err
	}

	log.Print("Start Listen Bot Message")
	go MessageReaderHolder(client.Done, client.msgGetChan, client.session, client.msgSocket, client.rootURL, &WSHolder{Conn: client.msgSocket})
	go MessageReaderHolder(client.Done, client.msgGetChan, client.session, client.msgSocket, client.rootURL, &WSHolder{Conn: client.msgSocket})
	log.Print("Bot Message Listener Started!")

	//TODO: Event Handle

	log.Print("Strat Lisen Message Send Channal")
	go MessageSenderHolder(client.Done, client.msgReq, client.session)
	go MessageSenderHolder(client.Done, client.msgReq, client.session)
	log.Print("Bot Message Send Channal Litener Started")

	return nil
}

func (client *BotWsClient) Close() {
	client.msgSocket.Close()
	client.eventSocket.Close()
}

func (client *BotWsClient) GetDoneChan() chan struct{} {
	return client.Done
}

func (client *BotWsClient) StopClient() {
	close(client.msgGetChan)
	close(client.msgReq)
	close(client.Done)

	log.Print("Release Session")
	releaseSessionBody := structs.VerifyQQRequest{
		QQ:         3628862306,
		SessionKey: string(client.session),
	}
	var res structs.VerifyRespond

	err := request.PostWithTargetRespond("/release", releaseSessionBody, &res)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("exit release session")

	client.msgSocket.WriteMessage(websocket.TextMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
}
