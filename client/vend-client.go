package main

import (
	"github.com/lexicality/vending/client/vendio"
	"github.com/lexicality/vending/shared"
)

var log = shared.GetLogger("client")

func main() {
	log.Info("Startup!")
	hardware := vendio.GetHardware(log)
	hardware.Setup()
	defer hardware.Teardown()
	wsHandler("ws://localhost:8080/ws", hardware)
}
