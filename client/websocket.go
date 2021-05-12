package client

import (
	"fmt"
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

type WebsocketDial func(url url.URL, session Session, msgSocket *websocket.Conn) error

func EstablishMessageHandleWebSocket(url url.URL, session Session, msgSocket *websocket.Conn) error {
	//check old init
	msgSocket.Close()

	s, err := establisWebSocket(url, Session(session), "/message")
	msgSocket = s
	return err
}

func EstablishEventHandleWebSocket(url url.URL, session Session, eventSocket *websocket.Conn) error {
	eventSocket.Close()

	s, err := establisWebSocket(url, Session(session), "/event")
	eventSocket = s
	return err
}

func establisWebSocket(url url.URL, session Session, path string) (*websocket.Conn, error) {
	targetURL := url
	targetURL.Path = "/message"
	targetURL.RawQuery = fmt.Sprintf("sessionKey=%s", session)

	s, _, err := websocket.DefaultDialer.Dial(targetURL.String(), nil)

	if err != nil {
		log.Fatal("Dial To websorct Falure Error: ", err)
		return s, err
	}
	return s, nil
}

func TryReDialWebSocket(dialer WebsocketDial, tryTime uint, session Session, conn *websocket.Conn, url url.URL) error {
	log.Printf("Trying to reConnect to websocket(0/%v)", tryTime)
	errs := []error{}
	for i := uint(0); i < tryTime; i++ {
		err := dialer(url, session, conn)
		if err != nil {
			log.Fatalf("Try ReConnect To Websoeck Failure(%v/%v): %v", tryTime, i+1, err.Error())
			errs = append(errs, err)
			continue
		}
		return nil
	}
	return fmt.Errorf("errors: %v", errs)
}
