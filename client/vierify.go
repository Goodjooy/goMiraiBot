package client

import (
	"fmt"
	"goMiraiQQBot/constdata"
	"goMiraiQQBot/request"
	"goMiraiQQBot/request/structs"
	"log"
)

type Session string

func AuthQQKey(cfg Config) (Session, error) {
	log.Println("send auth key")
	//get session
	var resInterface structs.AuthRespond
	err := request.PostWithTargetRespond(
		"/auth",
		request.H{
			"authKey": cfg.Bot.AuthKey,
		},
		&resInterface)

	if err != nil {
		log.Fatal("Send Auth Key Error: ", err)
		return "", err
	}
	if resInterface.Code != constdata.Normal {
		log.Fatal("Send Auth Key Error: StatusCode Error", resInterface.Code)
		return "", fmt.Errorf("ErrorCode: %v", resInterface.Code)
	}
	log.Println("Authentication Success")

	return Session(resInterface.Session), nil
}

func VerifyQQ(session Session, cfg Config) error {
	log.Print("Starting QQ Bot")

	log.Println("Verifying QQ")
	//verify qq
	verifyRequestBody := structs.VerifyQQRequest{
		QQ:         uint64(cfg.Bot.QQ),
		SessionKey: string(session),
	}
	var res structs.VerifyRespond
	err := request.PostWithTargetRespond("/verify", verifyRequestBody, &res)
	if err != nil || res.Code != constdata.Normal {
		log.Fatal("Verify QQ Bot Error: ", err)
		return err
	}
	log.Println("Verify QQ Success!")
	return nil
}
