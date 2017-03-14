package hardware

import (
	"context"
	"sync"

	logging "github.com/op/go-logging"

	"github.com/lexicality/vending/vend"
)

// Machine is a vending machine
type Machine struct {
	sync.Mutex

	hw  hardware
	log *logging.Logger
}

// NewMachine returns a new machine
func NewMachine(log *logging.Logger) *Machine {
	return &Machine{
		log: log,
	}
}

// SetupHardware configures the vending hardware and dispatches a worker thread to tear it down again when the context is canceled
func (m *Machine) SetupHardware(ctx context.Context) (err error) {
	m.hw, err = setupHardware(ctx, m.log)
	return err
}

// Vend is a temporary method that does the HW vend but lockishly
func (m *Machine) Vend(ctx context.Context, location uint8) vend.Result {
	m.Lock()
	defer m.Unlock()

	// Check if the session died while we were waiting
	err := ctx.Err()
	if err != nil {
		m.log.Criticalf("Aborting incomplete vend %s due to the context closing (%s)", "TODO", err)
		return vend.ResultAborted
	}

	return m.hw.Vend(ctx, location)
}
