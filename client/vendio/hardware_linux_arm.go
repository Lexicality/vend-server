package vendio

import (
	"time"

	"github.com/lexicality/vending/shared/vending"
	"github.com/op/go-logging"
	rpio "github.com/stianeikeland/go-rpio"
)

// Bottom 16 pins starting at physical pin 21
var outPins = [vending.MaxLocations]rpio.Pin{
	9, 25, 11, 8,
	7, 0, 1, 5,
	6, 12, 13, 19,
	16, 26, 20, 21,
}

// Two bit input status
const (
	inLow  rpio.Pin = 23 // 16
	inHigh rpio.Pin = 24 // 18
)

type MotorMode uint8

const (
	MotorOff MotorMode = iota
	MotorOn
	MotorJammed
	MotorEmpty
)

const (
	VendTime = time.Second * 30
	// VendCheckInterval = time.Millisecond * 500
	VendCheckInterval = time.Second
)

func sprintPinMode(mode rpio.State) string {
	if mode == rpio.High {
		return "High"
	}
	return "Low"
}

type hardware struct {
	log *logging.Logger
}

func (hw *hardware) Setup() error {
	if hw.log != nil {
		hw.log.Info("Hello I'm ARM!")
	}

	err := rpio.Open()
	if err != nil {
		return err
	}

	for _, pin := range outPins {
		pin.Output()
		pin.Low()
	}

	inHigh.Input()
	inLow.Input()

	return nil
}

func (hw *hardware) Teardown() error {
	return rpio.Close()
}

func (hw *hardware) getMotorMode() MotorMode {
	highPin := inHigh.Read()
	lowPin := inLow.Read()
	if hw.log != nil {
		hw.log.Debugf("MOTOR STATE PINS: %s %s", sprintPinMode(highPin), sprintPinMode(lowPin))
	}
	// TODO: Know something john snow
	return MotorOff

	h := highPin == rpio.High
	l := lowPin == rpio.Low
	// This is made up and is probably wrong
	if h && l {
		return MotorJammed
	} else if h && !l {
		return MotorOn
	} else if !h && l {
		return MotorEmpty
	} else {
		return MotorOff
	}
}

func (hw *hardware) Vend(location uint8) error {
	if hw.log != nil {
		hw.log.Infof("~~~I AM VENDING ITEM #%d!", location)
	}

	if location > vending.MaxLocations {
		return ErrInvalidLocation
	}

	hw.log.Debugf("INPUT HIGH BIT: %")

	for _, pin := range outPins {
		pin.Low()
	}
	outPins[location].High()
	// Always stop
	defer outPins[location].Low()

	// It takes ${VendTime} seconds to push out an item under normal circumstances
	endTimer := time.NewTimer(VendTime)
	defer endTimer.Stop()
	// Check every ${VendCheckInterval} that we've not become jammed
	checkTicker := time.NewTicker(VendCheckInterval)
	defer checkTicker.Stop()

	// Trigger the ticker check immediately
	t := make(chan bool)
	defer close(t)
	t <- true
	for {
		select {
		case <-endTimer.C:
			return nil
		case <-t:
		case <-checkTicker.C:
			motorState := hw.getMotorMode()
			if motorState == MotorJammed {
				return ErrMachineJammed
			} else if motorState == MotorEmpty {
				// TODO: If it shows up as empty after 29 seconds of not being empty it's probably a successful vend
				return ErrLocationEmpty
			}
		}
	}
}
