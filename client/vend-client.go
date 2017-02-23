package main

import (
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("Vend-client")

func main() {
	log.Info("Startup!")
	wsHandler("ws://localhost:8080/ws")
}
