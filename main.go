package main

import (
	"os"
	"os/signal"
	"syscall"

	"context"

	"github.com/lexicality/vending/backend"
	"github.com/lexicality/vending/hardware"
	"github.com/lexicality/vending/web"
)

const (
	// Development location of HTML etc etc
	webRoot = "src/github.com/lexicality/vending/web/www-src"
)

func ctrlCHandler(ctx context.Context) context.Context {
	newCtx, shutdown := context.WithCancel(context.Background())

	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		for {
			<-c
			if newCtx.Err() != nil {
				// Double ^c == GTFO
				os.Exit(1)
			} else {
				log.Info("Shutting down due to SIGINT")
				shutdown()
			}
		}
	}()

	return newCtx
}

func main() {
	ctx := ctrlCHandler(context.Background())
	setupLogging(ctx, "Vending")
	log.Info("Hello World")

	hw, err := hardware.SetupHardware(ctx, log)
	if err != nil {
		log.Fatalf("Unable to open vending hardware: %s", err)
	}

	stock := backend.GetFakeStock()
	web.Server(ctx, ":80", webRoot, log, stock, hw)
}
