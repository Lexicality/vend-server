package main

import (
	"github.com/op/go-logging"
)

var log *logging.Logger

func setupLogging(name string) *logging.Logger {
	log = logging.MustGetLogger(name)
	// TODO: Set up logging formats
	return log
}
