package main

import (
	"encoding/json"

	"github.com/gorilla/websocket"
	"github.com/lexicality/vending/shared"
	"github.com/lexicality/vending/shared/vending"
)

func handleMessage(message []byte) error {
	msg := vending.RecvMessage{}
	err := json.Unmarshal(message, &msg)
	if err != nil {
		return err
	} else if msg.Type != "Request" {
		log.Warningf("Unahandled message %s with type %s!", msg.Message, msg.Type)
		return nil
	}

	req := vending.Request{}
	err = json.Unmarshal(msg.Message, &req)
	if err != nil {
		return err
	}

	vendItem(req.Location)

	return nil
}

func readPump(conn *shared.WSConn) error {
	for {
		mType, msg, err := conn.ReadMessage()
		if err != nil {
			return err
		} else if mType != websocket.TextMessage {
			log.Warningf("Got unknown message type %+v with message %s", mType, msg)
			continue
		}

		conn.MessageRecieved()
		log.Debugf("MESSAGE: %s", msg)

		err = handleMessage(msg)
		if err != nil {
			log.Warningf("Unable to handle message %s: %s", msg, err)
			continue
		}
	}
}

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

	err = conn.WriteJSON(&vending.SendMessage{
		Type:    "Welcome",
		Message: "Hello!",
	})
	if err != nil {
		log.Fatalf("It's not actually open :(")
	}

	err = readPump(conn)
	if err != nil {
		log.Errorf("Connection died: %s", err)
	} else {
		log.Error("Connection died?")
	}
}
