package hardware

import (
	"context"
	"sync"

	logging "github.com/op/go-logging"

	"github.com/lexicality/vending/vend"
)

// Machine is a cross-platform way of using the vending machine
type Machine interface {
	// The Machine will be locked whenever vending, exposed for when actions need to wait for the current vend.
	sync.Locker
	// Setup configures the machine to do things, and sets up a handler to stop doing things when the context closes.
	// If it returns an error all following vends will fail.
	Setup(ctx context.Context) error
	// Vend requests the hardware to vend an item. Blocks until done.
	// If a second request comes in while the first is processing that will block too.
	// Vend respects context deadlines / cancels before the physical vend action starts,
	//  but will not abort once in progress to avoid becoming physically out of sync.
	Vend(ctx context.Context, location uint8) vend.Result
}

// NewMachine returns this platform's vending options
func NewMachine(log *logging.Logger) Machine {
	return &physicalHardware{
		log: log,
	}
}
