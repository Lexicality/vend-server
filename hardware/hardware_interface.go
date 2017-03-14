package hardware

import (
	"context"

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
	// If the passed context has a deadline sooner than the expected vend time it aborts immediately. (Cannot be canceled)
	Vend(ctx context.Context, location uint8) vend.Result
}

func hwmonitor(ctx context.Context, log *logging.Logger, hw Hardware) {
	<-ctx.Done()
	err := hw.Teardown()
	if err != nil && log != nil {
		log.Criticalf("Unable to free HW data: %s", err)
	}
}

// SetupHardware configures and returns a Hardware instance
func SetupHardware(ctx context.Context, log *logging.Logger) (Hardware, error) {
	hw := &physicalHardware{
		log: log,
	}

	err := hw.Setup()
	if err != nil {
		return nil, err
	}
	go hwmonitor(ctx, log, hw)

	return hw, nil
}
