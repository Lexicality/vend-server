package main

import (
	"fmt"

	"github.com/go-martini/martini"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("test")

func main() {
	fmt.Println("Hello World")
	m := martini.Classic()
	m.Get("/", func() string {
		return "Hello world!"
	})
	m.Get("/ws", wsHandler)
	m.RunOnAddr(":8080")
}
