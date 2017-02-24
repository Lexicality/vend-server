package shared

import (
	"github.com/op/go-logging"
)

var log *logging.Logger

// GetLogger gets a logger in the state what I like it in.
// Also it sets the global for the shared package
func GetLogger(name string) *logging.Logger {
	log = logging.MustGetLogger(name)
	// TODO: Set up logging formats
	return log
}
