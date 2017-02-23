package main

import (
	"github.com/gorilla/websocket"
)

var dialer = websocket.Dialer{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Discard function to keep socket alive
func readLoop(c *websocket.Conn) {
	for {
		if _, _, err := c.NextReader(); err != nil {
			log.Noticef("Connection closed: %s", err)
			c.Close()
			break
		}
	}
}

func wsHandler(server string) {
	log.Notice("Connection attempt begining")
	var err error

	conn, _, err := dialer.Dial(server, nil)
	if err != nil {
		log.Errorf("Unable to connect to server: %s", err)
		return
	}
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

		log.Debugf("MESSAGE: %s", msg)
	}
}
