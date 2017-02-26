package main

import (
	"github.com/gorilla/websocket"
	"github.com/lexicality/vending/shared"
)

func wsHandler(server string) {
	log.Notice("Connection attempt begining")
	var err error

	c, _, err := websocket.DefaultDialer.Dial(server, nil)
	if err != nil {
		log.Errorf("Unable to connect to server: %s", err)
		return
	}
	var conn = shared.NewWSConn(c)
	defer conn.Close()

	err = conn.WriteMessage(websocket.TextMessage, []byte("hi!"))
	if err != nil {
		log.Fatalf("It's not actually open :(")
	}

	for {
		mType, msg, err := conn.ReadMessage()
		if err != nil {
			log.Errorf("Unable to read next message: %s", err)
			break
		} else if mType != websocket.TextMessage {
			log.Warningf("Got unknown message type %+v with message %s", mType, msg)
			continue
		}

		conn.MessageRecieved()
		log.Debugf("MESSAGE: %s", msg)
	}
}
