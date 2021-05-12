package main

import (
	"flag"
	"goMiraiQQBot/client"
	"goMiraiQQBot/messageHandle/interact"
	"goMiraiQQBot/messageHandle/interactHandle"
	appchaininteract "goMiraiQQBot/messageHandle/interactHandle/appChainInteract"
	hentaiimageinteract "goMiraiQQBot/messageHandle/interactHandle/hentaiImageInteract"
	xmlchaininteract "goMiraiQQBot/messageHandle/interactHandle/xmlChainInteract"
	"log"
	"os"
	"os/signal"
	"time"
)

func main() {

	flag.Parse()
	log.SetFlags(1)

	//close workser
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	client := client.NewQQBotClient("./config.yml")

	defer client.StopClient()
	defer client.Close()

	interact.AddSingleInteract(interactHandle.NewHelpInteract)
	interact.AddSingleInteract(interactHandle.NewAboutInteract)

	interact.AddContextInteract(hentaiimageinteract.NewHentaiImageSearchInteract)

	interact.AddChainSingleInteract(xmlchaininteract.NewXmlChainInteract)
	interact.AddChainSingleInteract(appchaininteract.NewAppChainInteract)

	//load message
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
	END:return
}
