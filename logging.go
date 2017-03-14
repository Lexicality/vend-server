package main

import (
	"context"

	logging "github.com/op/go-logging"
)

var log *logging.Logger

func setupLogging(ctx context.Context, name string) *logging.Logger {
	log = logging.MustGetLogger(name)
	// TODO: Set up logging formats w/ syslog etc
	return log
}
