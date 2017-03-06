package main

import (
	"fmt"

	"github.com/lexicality/vending/shared"
)

var log = shared.GetLogger("server")

const (
	// Development location of HTML etc etc
	webRoot = "src/github.com/lexicality/vending/server/www-src"
)

func main() {
	fmt.Println("Hello World")

	stock := GetFakeStock()
	// stdinStream := readStreamer()

	messagePub(stock.VendC)
	go handlePubSub()
	go wsServer(":8080")
	webServer(":80", webRoot, stock)
}
