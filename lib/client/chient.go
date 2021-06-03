package client

import (
	"fmt"
	"goMiraiQQBot/lib/database"
	ih "goMiraiQQBot/lib/interactHandle"
	appchaininteract "goMiraiQQBot/lib/interactHandle/appChainInteract"
	xmlchaininteract "goMiraiQQBot/lib/interactHandle/xmlChainInteract"
	"goMiraiQQBot/lib/messageHandle/interact"
	messagetargets "goMiraiQQBot/lib/messageHandle/messageTargets"
	"goMiraiQQBot/lib/messageHandle/sourceHandle"
	mStruct "goMiraiQQBot/lib/messageHandle/structs"
	"goMiraiQQBot/lib/request"
	"goMiraiQQBot/lib/request/structs"
	"log"
	"net/url"

	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

type BotWsClient struct {
	rootURL url.URL

	msgSocket   *websocket.Conn
	eventSocket *websocket.Conn

	session Session

	Done chan struct{}

	msgChan MessageChan

	config Config

	database *gorm.DB

	CmdBaseInteract interact.InteractController
	ChainInteract interact.InteractController
}

func NewQQBotClient(configPath string) BotWsClient {
	client := BotWsClient{
		Done:       make(chan struct{}),
		msgChan: MessageChan{
			inputMsg: make(chan mStruct.Message),
			outputMsg: make(chan messagetargets.MessageTarget),
		},
		
		CmdBaseInteract: interact.NewCommandBaseInteractGroup(),
		ChainInteract: interact.ChainBaseInteractController(),
	}

	config := LoadConfig(configPath)
	client.config = config
	cfg = config
	interact.SetCFG(config)

	//TODO: 控制内置
	//client.CmdBaseInteract.AddSingleInteractConstruct(ih.NewHelpInteract)
	interact.MessageInteract.AddSingleConstruct(ih.NewHelpInteract)
	interact.MessageInteract.AddSingleConstruct(ih.NewAboutInteract)

	interact.ChainInteract.AddSingleConstruct(xmlchaininteract.NewXmlChainInteract)
	interact.ChainInteract.AddSingleConstruct(appchaininteract.NewAppChainInteract)

	return client
}

func (client *BotWsClient) EnableDB() {
	database.Init(client.config)
}

func (client *BotWsClient) Run() error {
	log.Print("Starting QQ Bot")
	config := client.config

	log.Print("Init Souce Handle...")
	sourceHandle.InitMessageSourceHandle()
	log.Print("Souce Handle Init Success!")

	log.Print("Init Interact Handle")
	interact.InitInteractHandle(client.msgChan.inputMsg, client.msgChan.outputMsg)
	log.Print("Interact Handle Init Success")

	session, err := AuthQQKey(config)
	if err != nil {
		log.Print("Init Bot Clent Faulure | Send Auth Key: ", err)
		return err
	} else {
		client.session = session
	}
	err = VerifyQQ(session, config)
	if err != nil {
		log.Print("Init Bot Clent Faulure | Verify QQ: ", err)
		return err
	}

	log.Print("Establish Websocket Connect To Bot")
	client.rootURL = url.URL{Host: fmt.Sprintf("%v:%v", config.Server.Host, config.Server.Port)}

	client.msgSocket, err = EstablishMessageHandleWebSocket(client.rootURL, session)
	if err != nil {
		log.Print("Init Bot Clent Faulure | Establish Websocket at `/message` : ", err)
		return err
	}

	log.Print("Start Listen Bot Message")
	go MessageReaderHolder(client.Done, client.msgChan.inputMsg, client.session, client.msgSocket, client.rootURL, &WSHolder{Conn: client.msgSocket})
	log.Print("Bot Message Listener Started!")

	//TODO: Event Handle

	log.Print("Strat Lisen Message Send Channal")
	go MessageSenderHolder(client.Done, client.msgChan.outputMsg, client.session)
	go MessageSenderHolder(client.Done, client.msgChan.outputMsg, client.session)
	log.Print("Bot Message Send Channal Litener Started")

	return nil
}

func (client *BotWsClient) Close() {
	client.msgSocket.Close()
	//client.eventSocket.Close()
}

func (client *BotWsClient) GetDoneChan() chan struct{} {
	return client.Done
}

func (client *BotWsClient) StopClient() {
	close(client.msgChan.inputMsg)
	close(client.msgChan.outputMsg)
	close(client.Done)

	log.Print("Release Session")
	releaseSessionBody := structs.VerifyQQRequest{
		QQ:         3628862306,
		SessionKey: string(client.session),
	}
	var res structs.VerifyRespond

	err := request.PostWithTargetRespond("/release", releaseSessionBody, &res)
	if err != nil {
		log.Fatalf("send Message Fauilure: %v", err)
	}
	log.Print("exit release session")

	client.msgSocket.WriteMessage(websocket.TextMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
}
