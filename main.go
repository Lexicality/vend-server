package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	stripe "github.com/stripe/stripe-go"

	"github.com/lexicality/vending/backend"
	"github.com/lexicality/vending/hardware"
	"github.com/lexicality/vending/web"
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
	var err error
	ctx := ctrlCHandler(context.Background())
	setupLogging(ctx, "Vending")
	log.Info("Hello World")

	hw := hardware.NewMachine(log)
	err = hw.SetupHardware(ctx)
	if err != nil {
		log.Fatalf("Unable to open vending hardware: %s", err)
	}

	stock := backend.GetFakeStock()

	go web.ServeCanonical(":http", "https://vend.lan.london.hackspace.org.uk")

	stripe.Key = "sk_test_aB3qhPiZpdnv31WnNmI8DFFR"

	webServer := &web.Server{
		// TODO: All of this should be configured by flags / env vars / whatever
		Addr:        ":https",
		WebRoot:     "src/github.com/lexicality/vending/web/www-src",
		ServerName:  "vend.lan.london.hackspace.org.uk",
		TLSCertFile: "cert.pem",
		TLSKeyFile:  "key.pem",
	}
	err = webServer.ServeHTTP(ctx, log, stock, hw)

	if err == http.ErrServerClosed {
		log.Infof("HTTP server shut down")
	} else if err == context.DeadlineExceeded {
		log.Criticalf("Web server timed out shutting down!")
	} else if err != nil {
		log.Fatalf("Web serving error: %s", err)
	}

	// Wait for vending machine to finish vending
	log.Noticef("Waiting for hardware shutdown")
	hw.Lock()
	hw.Unlock()

	// Wait for any pending vends that may exist to die
	log.Noticef("Waiting for context drain")
	<-time.After(time.Millisecond * 100)

	log.Debug("(✖╭╮✖)")
}
