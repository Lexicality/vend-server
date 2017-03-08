package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/lexicality/vending/client/vendio"
	"github.com/lexicality/vending/shared"
	"github.com/lexicality/vending/shared/vending"
)

func handleMessage(hw vendio.Hardware, msg *vending.RecvMessage) (resp *vending.SendMessage, err error) {
	switch msg.Type {
	case "Request":
		req := vending.Request{}
		err = json.Unmarshal(msg.Message, &req)
		if err != nil {
			return nil, err
		}

		log.Infof("Vending location %d for request %s", req.Location, req.ID)

		res := hw.Vend(req.Location)

		return &vending.SendMessage{
			Type: "Response",
			Message: &vending.Response{
				ID:     req.ID,
				Result: res,
			},
		}, nil
	default:
		log.Warningf("Unahandled message %s with type %s!", msg.Message, msg.Type)
		return nil, nil
	}
}

func readPump(conn *shared.WSConn, hw vendio.Hardware) error {
	var err error
	// Reuse the same message object
	var msg = &vending.RecvMessage{}
	var resp *vending.SendMessage
	for {
		err = conn.ReadJSON(msg)
		if err != nil {
			return err
		}
		conn.MessageRecieved()
		log.Debugf("Recieved %s message: %s", msg.Type, msg.Message)

		resp, err = handleMessage(hw, msg)
		if err != nil {
			log.Errorf("Unable to handle message %s: %s", msg, err)
			continue
		} else if resp != nil {
			err = conn.WriteJSON(resp)
			if err != nil {
				return err
			}
		}
	}
}

var wsDialer = websocket.Dialer{
	Proxy:        http.ProxyFromEnvironment,
	Subprotocols: []string{vending.MessageProtocol},
}

func wsHandler(server string, hw vendio.Hardware) {
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

	err = readPump(conn, hw)
	if err != nil {
		log.Errorf("Connection died: %s", err)
	} else {
		log.Error("Connection died?")
	}
}
