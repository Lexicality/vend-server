package main

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/lexicality/vending/shared"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Disable security until I sort out the client
	CheckOrigin: func(r *http.Request) bool { return true },
}

func wsWriteLoop(conn *shared.WSConn) {

	msgChan := messageSub()
	pingChang := conn.GetPingTicker()

	var msg string
	var err error
	var ok bool
	for {
		select {
		case msg, ok = <-msgChan:
			if !ok {
				log.Info("Killing connection due to channel closure")
				return
			} else if !conn.IsOpen() {
				log.Debug("Killing write loop due to connection closure")
				return
			}

			err = conn.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				log.Error("Unable to send message %s: %s", msg, err)
				return
			}
		case _ = <-pingChang:
			if conn.IsTimingOut() {
				log.Info("Killing connection due to timeout")
				return
			} else if !conn.IsOpen() {
				log.Debug("Killing write loop due to connection closure")
				return
			}

			err = conn.MaybeSendPing()
			if err != nil {
				log.Error("Unable to send ping: %s", err)
				return
			}
		}
	}
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	log.Notice("Connection attempt begining")
	var err error

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error(err)
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}
	var conn = shared.NewWSConn(c)
	defer c.Close()

	_ = conn.WriteMessage(websocket.TextMessage, []byte("hi!"))

	// Ignore anything the client has to say
	go conn.ReadDiscardPump()
	// Tell them all the important things we have to say
	wsWriteLoop(conn)
}
