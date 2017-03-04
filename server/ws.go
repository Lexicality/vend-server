package main

import (
	"net/http"

	"net"

	"github.com/gorilla/websocket"
	"github.com/lexicality/vending/shared"
	"github.com/lexicality/vending/shared/vending"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	Subprotocols:    []string{vending.MessageProtocol},
	// Disable security until I sort out the client
	CheckOrigin: func(r *http.Request) bool { return true },
}

func wsReadLoop(conn *shared.WSConn) {
	var err error
	var msg = &vending.RecvMessage{}
	for {
		err = conn.ReadJSON(msg)
		if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
			log.Debug("Ending read loop due to socket closing")
			return
		} else if err != nil {
			log.Infof("Closing connection due to read error: %s", err)
			conn.Close()
			return
		}

		conn.MessageRecieved()
		log.Debugf("Recieved %s message: %s", msg.Type, msg.Message)

		// TODO
	}
}

func wsWriteLoop(conn *shared.WSConn) {

	msgChan := messageSub()
	pingChang := conn.GetPingTicker()

	var msg *vending.SendMessage
	var err error
	var ok bool
	for {
		err = nil
		msg = nil
		select {
		case msg, ok = <-msgChan:
			if !ok {
				log.Info("Killing connection due to channel closure")
				conn.Close()
				return
			} else if !conn.IsOpen() {
				log.Debug("Killing write loop due to connection closure")
				return
			}

			conn.SetWriteDeadline(conn.GetWriteDeadline())
			err = conn.WriteJSON(msg)
		case _ = <-pingChang:
			if conn.IsTimingOut() {
				log.Info("Killing connection due to ping timeout")
				conn.Close()
				return
			} else if !conn.IsOpen() {
				log.Debug("Killing write loop due to connection closure")
				return
			}

			err = conn.MaybeSendPing()
		}

		if err == nil {
			continue
		}

		// ERROR HANDLING YEAH
		switch v := err.(type) {
		case net.Error:
			if v.Timeout() {
				log.Info("Killing connection due to write timeout")
				conn.Close()
				return
			}
		case *websocket.CloseError:
			// ?
		}

		// "idk lol"
		if msg != nil {
			log.Errorf("Unable to send message %+v: %s", msg, err)
		} else {
			log.Errorf("Unable to send message: %s", err)
		}
	}
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	log.Notice("Connection attempt begining")
	var err error

	subProtocols := websocket.Subprotocols(r)
	if len(subProtocols) != 1 || subProtocols[0] != vending.MessageProtocol {
		log.Noticef("Connection attempt with wrong protocol - expected [%s] but got %v", vending.MessageProtocol, subProtocols)
		w.WriteHeader(400)
		w.Write([]byte("Bad protocol"))
		return
	}

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error(err)
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}
	var conn = shared.NewWSConn(c)
	defer c.Close()

	// Make sure the socket is open
	err = conn.WriteJSON(&vending.SendMessage{
		Type:    "Welcome",
		Message: "Hello!",
	})
	if err != nil {
		log.Warningf("Unable to send welcome message: %s", err)
		return
	}

	// Ignore anything the client has to say
	go wsReadLoop(conn)
	// Tell them all the important things we have to say
	wsWriteLoop(conn)
}
