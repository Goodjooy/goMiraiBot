package main

import (
	"flag"
	"goMiraiQQBot/client"
	"goMiraiQQBot/interactHandle"
	appchaininteract "goMiraiQQBot/interactHandle/appChainInteract"
	hentaiimageinteract "goMiraiQQBot/interactHandle/hentaiImageInteract"
	inandoutinteract "goMiraiQQBot/interactHandle/inAndOutInteract"
	xmlchaininteract "goMiraiQQBot/interactHandle/xmlChainInteract"
	"goMiraiQQBot/messageHandle/interact"
	"log"
	"os"
	"os/signal"
	"time"
)

func main() {

	flag.Parse()
	log.SetFlags(0)

	//close workser
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	client := client.NewQQBotClient("./config.yml")

	defer client.StopClient()
	defer client.Close()

	interact.MessageInteract.AddSingleConstruct(interactHandle.NewHelpInteract)
	interact.MessageInteract.AddSingleConstruct(interactHandle.NewAboutInteract)

	interact.MessageInteract.AddContextConstruct(hentaiimageinteract.NewHentaiImageSearchInteract)
	interact.MessageInteract.AddContextConstruct(inandoutinteract.NewPaymentRecordInteract)

	interact.ChainInteract.AddSingleConstruct(xmlchaininteract.NewXmlChainInteract)
	interact.ChainInteract.AddSingleConstruct(appchaininteract.NewAppChainInteract)

	//load message
	client.EnableDB()
	err := client.Run()
	if err != nil {
		log.Fatal("Run Client Failure: ", err)
		return
	}

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-client.Done:
			goto END
		case <-interrupt:
			{
				log.Print("interrupt")
				goto END
			}
		}
	}
END:
	return
}
