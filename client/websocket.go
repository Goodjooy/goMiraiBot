package client

import (
	"fmt"
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

type WebsocketDial func(url url.URL, session Session) (*websocket.Conn, error)

func EstablishMessageHandleWebSocket(url url.URL, session Session) (*websocket.Conn, error) {
	//check old init

	s, err := establisWebSocket(url, Session(session), "/message")
	return s, err
}

func EstablishEventHandleWebSocket(url url.URL, session Session) (*websocket.Conn, error) {

	s, err := establisWebSocket(url, Session(session), "/event")
	return s, err
}

func establisWebSocket(url url.URL, session Session, path string) (*websocket.Conn, error) {
	targetURL := url
	targetURL.Scheme = "ws"
	targetURL.Path = "/message"
	targetURL.RawQuery = fmt.Sprintf("sessionKey=%s", session)

	s, _, err := websocket.DefaultDialer.Dial(targetURL.String(), nil)

	if err != nil {
		log.Fatal("Dial To websorct Falure Error: ", err)
		return s, err
	}
	return s, nil
}

func TryReDialWebSocket(dialer WebsocketDial, tryTime uint, session Session, url url.URL) (*websocket.Conn, error) {
	log.Printf("Trying to reConnect to websocket(0/%v)", tryTime)
	errs := []error{}
	for i := uint(0); i < tryTime; i++ {
		conn, err := dialer(url, session)
		if err != nil {
			log.Fatalf("Try ReConnect To Websoeck Failure(%v/%v): %v", tryTime, i+1, err.Error())
			errs = append(errs, err)
			continue
		}
		return conn, nil
	}
	return &websocket.Conn{}, fmt.Errorf("errors: %v", errs)
}
