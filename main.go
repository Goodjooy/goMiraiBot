package main

import (
	"flag"
	"goMiraiQQBot/interactHandle/billi/upfans"
	hentaiimageinteract "goMiraiQQBot/interactHandle/hentaiImageInteract"
	inandoutinteract "goMiraiQQBot/interactHandle/inAndOutInteract"
	randomsetuinteract "goMiraiQQBot/interactHandle/randomSetuInteract"
	todaynewsinteractgo "goMiraiQQBot/interactHandle/todayNewsInteract.go"
	"goMiraiQQBot/lib/client"
	"goMiraiQQBot/lib/messageHandle/interact"
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

	interact.MessageInteract.AddSingleConstruct(randomsetuinteract.NewSetuInteract)
	interact.MessageInteract.AddSingleConstruct(todaynewsinteractgo.NewTodyNewsInteract)
	interact.MessageInteract.AddSingleConstruct(upfans.NewUpFansInteract)

	interact.MessageInteract.AddContextConstruct(hentaiimageinteract.NewHentaiImageSearchInteract)
	interact.MessageInteract.AddContextConstruct(inandoutinteract.NewPaymentRecordInteract)

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
