package main

import (
	"fmt"
	"net/http"

	"github.com/go-martini/martini"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("test")

func main() {
	fmt.Println("Hello World")

	messagePub(readStreamer())
	go handlePubSub()
	go tcpServer(":8081")

	m := martini.Classic()
	m.Get("/", func() string {
		return "Hello world!"
	})
	m.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "src/github.com/lexicality/vending/server/www-src/test.html")
	})
	m.Get("/ws", wsHandler)
	m.RunOnAddr(":8080")
	log.Fatal("Web server stopped?")
}
