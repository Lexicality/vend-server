package hardware

import (
	"context"
	"time"

	logging "github.com/op/go-logging"
	rpio "github.com/stianeikeland/go-rpio"

	"sync"

	"github.com/lexicality/vending/vend"
)

// Bottom 16 pins starting at physical pin 21
var outPins = [vend.MaxLocations]rpio.Pin{
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

type physicalHardware struct {
	sync.Mutex
	log *logging.Logger
}

func (hw *physicalHardware) Setup(ctx context.Context) error {
	if hw.log != nil {
		hw.log.Info("Hello I'm ARM!")
	}

	err := rpio.Open()
	if err != nil {
		return err
	}
	go hw.deferredTeardown(ctx)

	for _, pin := range outPins {
		pin.Output()
		pin.Low()
	}

	inHigh.Input()
	inLow.Input()

	return nil
}

func (hw *physicalHardware) deferredTeardown(ctx context.Context) {
	<-ctx.Done()
	err := rpio.Close()
	if err != nil && hw.log != nil {
		hw.log.Criticalf("Unable to free HW data: %s", err)
	}
}

func (hw *physicalHardware) getMotorMode() MotorMode {
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

func (hw *physicalHardware) Vend(ctx context.Context, location uint8) vend.Result {
	hw.Lock()
	defer hw.Unlock()

	err := ctx.Err()
	if err != nil {
		if hw.log != nil {
			hw.log.Warningf("Canceling vend attempt: %s", err)
		}
		return vend.ResultAborted
	} else if location > vend.MaxLocations {
		return vend.ResultInvalidRequest
	} else if dl, ok := ctx.Deadline(); ok && dl.Before(time.Now().Add(VendTime)) {
		// Don't even try and do anything if we're going to be aborted before we can vend
		if hw.log != nil {
			hw.log.Warningf("Canceling vend attempt as it would have timed out the context")
		}
		return vend.ResultAborted
	}

	if hw.log != nil {
		hw.log.Infof("~~~I AM VENDING ITEM #%d!", location)
	}

	// Dump debugging info before starting
	_ = hw.getMotorMode()

	// TODO: If the motor state is MotorOn, something has gone very wrong

	// TODO: This should not be necessary
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
		// We do *NOT* listen to context cancelations to avoid partial vends
		case <-endTimer.C:
			return vend.ResultSuccess
		case <-t:
		case <-checkTicker.C:
			motorState := hw.getMotorMode()
			if motorState == MotorJammed {
				return vend.ResultJammed
			} else if motorState == MotorEmpty {
				// TODO: If it shows up as empty after 29 seconds of not being empty it's probably a successful vend
				return vend.ResultEmpty
			}
			// TODO: If the motor state doesn't change to MotorOn then there's a hardware failure somewhere
		}
	}
}
