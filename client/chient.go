package client

import (
	"fmt"
	"goMiraiQQBot/messageHandle/interact"
	messagetargets "goMiraiQQBot/messageHandle/messageTargets"
	"goMiraiQQBot/messageHandle/sourceHandle"
	"goMiraiQQBot/messageHandle/structs"
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

	msgGetChan chan structs.Message
	msgReq     chan messagetargets.MessageTarget

	config Config
}

func NewQQBotClient(configPath string) BotWsClient {
	client := BotWsClient{
		Done:       make(chan struct{}),
		msgGetChan: make(chan structs.Message),
		msgReq:     make(chan messagetargets.MessageTarget),
	}

	config := LoadConfig(configPath)
	client.config = config

	return client
}

func (client *BotWsClient) Run(configPath string) error {
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

	err = EstablishMessageHandleWebSocket(client.rootURL, session, client.msgSocket)
	if err != nil {
		log.Fatal("Init Bot Clent Faulure | Establish Websocket at `/message` : ", err)
		return err
	}

	log.Print("Start Listen Bot Message")
	go MessageReaderHolder(client.Done, client.msgGetChan, client.session, client.msgSocket, client.rootURL)
	go MessageReaderHolder(client.Done, client.msgGetChan, client.session, client.msgSocket, client.rootURL)
	log.Print("Bot Message Listener Started!")

	//TODO: Event Handle

	log.Print("Strat Lisen Message Send Channal")
	go MessageSenderHolder(client.Done, client.msgReq, client.session)
	go MessageSenderHolder(client.Done, client.msgReq, client.session)
	log.Print("Bot Message Send Channal Litener Started")

	return nil
}
