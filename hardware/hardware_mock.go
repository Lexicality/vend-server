package hardware

import (
	"context"
	"sync"

	"github.com/lexicality/vending/vend"
)

// MockHardware is a mock for the Hardware type
type MockHardware struct {
	sync.Mutex
	SetupError  error
	VendResult  vend.Result
	VendRequest uint8
}

func (hw *MockHardware) Setup(ctx context.Context) error {
	return hw.SetupError
}

func (hw *MockHardware) Vend(ctx context.Context, location uint8) vend.Result {
	// Doesn't seem much point in doing this but it's in the spec
	hw.Lock()
	defer hw.Unlock()

	hw.VendRequest = location

	return hw.VendResult
}
