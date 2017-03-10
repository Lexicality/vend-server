package hardware

import (
	"github.com/lexicality/vending/vend"
	"github.com/op/go-logging"
)

// Hardware represents the actual vending IO interface
type Hardware interface {
	// Setup prepares the GPIO pins etc
	Setup() error
	// Teardown closes anything required to set up the GPIO
	Teardown() error
	// Vend requests the hardware to vend an item. Blocks until done.
	Vend(location uint8) vend.Result
}

// GetHardware returns an appropriate Hardware for this system
func GetHardware(log *logging.Logger) Hardware {
	return &hardware{
		log: log,
	}
}
