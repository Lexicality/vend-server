package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func readStreamer() <-chan string {
	stream := make(chan string)
	reader := bufio.NewReader(os.Stdin)
	go func() {
		for {
			fmt.Print("Enter message to send: ")
			text, err := reader.ReadString('\n')
			if err == io.EOF {
				log.Fatal("DONE")
			} else if err != nil {
				log.Fatal("Unable to read from stdin: %s", err)
			}

			text = strings.TrimSpace(text)

			// Never send an empty string since that's what closed channels do
			if text != "" {
				stream <- text
			}
		}
	}()

	return stream
}
