package main

import (
	"fmt"

	"github.com/lexicality/vending/backend"
	"github.com/lexicality/vending/hardware"
	"github.com/lexicality/vending/web"
)

const (
	// Development location of HTML etc etc
	webRoot = "src/github.com/lexicality/vending/web/www-src"
)

func main() {
	setupLogging("Vending")
	fmt.Println("Hello World")

	hw := hardware.GetHardware(log)
	err := hw.Setup()
	if err != nil {
		log.Fatalf("Unable to open vending hardware: %s", err)
	}
	// TODO: This can error
	defer hw.Teardown()

	stock := backend.GetFakeStock()
	web.Server(":80", webRoot, log, stock, hw)
}
