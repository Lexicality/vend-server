package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/lexicality/vending/shared"
	"github.com/lexicality/vending/shared/vending"
)

func handleMessage(msg *vending.RecvMessage) error {
	if msg.Type != "Request" {
		log.Warningf("Unahandled message %s with type %s!", msg.Message, msg.Type)
		return nil
	}

	req := vending.Request{}
	err := json.Unmarshal(msg.Message, &req)
	if err != nil {
		return err
	}

	vendItem(req.Location)

	return nil
}

func readPump(conn *shared.WSConn) error {
	var err error
	// Reuse the same message object
	var msg = &vending.RecvMessage{}
	for {
		err = conn.ReadJSON(msg)
		if err != nil {
			return err
		}
		conn.MessageRecieved()
		log.Debugf("Recieved %s message: %s", msg.Type, msg.Message)

		err = handleMessage(msg)
		if err != nil {
			log.Warningf("Unable to handle message %s: %s", msg, err)
			continue
		}
	}
}

var wsDialer = websocket.Dialer{
	Proxy:        http.ProxyFromEnvironment,
	Subprotocols: []string{vending.MessageProtocol},
}

func wsHandler(server string) {
	log.Notice("Connection attempt begining")
	var err error

	c, _, err := wsDialer.Dial(server, nil)
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
