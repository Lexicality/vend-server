package hardware

import (
	"github.com/lexicality/vending/vend"
)

// MockHardware is a mock for the Hardware type
type MockHardware struct {
	SetupError, TeardownError error
	VendResult                vend.Result
	VendRequest               uint8
}

func (hw *MockHardware) Setup() error {
	return hw.SetupError
}

func (hw *MockHardware) Teardown() error {
	return hw.TeardownError
}

func (hw *MockHardware) Vend(location uint8) vend.Result {
	hw.VendRequest = location
	return hw.VendResult
}
