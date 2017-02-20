package main

import (
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
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

func wsHandler(w http.ResponseWriter, r *http.Request) {
	log.Notice("Connection attempt begining")
	var err error

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error(err)
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}
	defer conn.Close()

	go readLoop(conn)
	_ = conn.WriteMessage(websocket.TextMessage, []byte("hi!"))

	data := messageSub()

	var msg string
	for {
		msg = <-data

		// Check for being booted off the channel
		if msg == "" {
			log.Info("Killing connection due to channel closure")
			break
		}

		err = conn.WriteMessage(websocket.TextMessage, []byte(msg))
		if err != nil {
			log.Error("Unable to send message %s!", msg)
			return
		}
	}
}
