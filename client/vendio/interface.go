package vendio

import (
	"errors"

	"github.com/op/go-logging"
)

var (
	// ErrInvalidLocation is if you ask to vend something that's not there (eg #225 or #16)
	ErrInvalidLocation = errors.New("Invalid Vend Location")
	// ErrMachineJammed is if the jam detector triggers during a vend operation
	ErrMachineJammed = errors.New("Jam detected")
	// ErrLocationEmpty is if the empty detector triggers during a vend operation
	ErrLocationEmpty = errors.New("Location Empty")
)

// Hardware represents the actual vending IO interface
type Hardware interface {
	// Setup prepares the GPIO pins etc
	Setup() error
	// Teardown closes anything required to set up the GPIO
	Teardown() error
	// Vend requests the hardware to vend an item. Blocks until done.
	Vend(location uint8) error
}

// GetHardware returns an appropriate Hardware for this system
func GetHardware(log *logging.Logger) Hardware {
	return &hardware{
		log: log,
	}
}
