package vendio

// MockHardware is a mock for the Hardware type
type MockHardware struct {
	SetupError, TeardownError, VendError error
	VendRequest                          uint8
}

func (hw *MockHardware) Setup() error {
	return hw.SetupError
}

func (hw *MockHardware) Teardown() error {
	return hw.TeardownError
}

func (hw *MockHardware) Vend(location uint8) error {
	hw.VendRequest = location
	return hw.VendError
}
