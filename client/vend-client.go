package main

import (
	"github.com/lexicality/vending/shared"
)

var log = shared.GetLogger("client")

func main() {
	log.Info("Startup!")
	wsHandler("ws://localhost:8080/ws")
}
