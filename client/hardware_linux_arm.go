package main

import (
	"github.com/stianeikeland/go-rpio"
)

// The Pi has a silly pin addressing system
const (
	pinStrobe rpio.Pin = 25 // 22
	pinData   rpio.Pin = 8  // 24
	pinClock  rpio.Pin = 7  // 26
	pinOE     rpio.Pin = 1  // 28
)

func readyHardware() (err error) {

	log.Info("Hello I'm ARM!")
	err = rpio.Open()
	if err != nil {
		return
	}
	pinStrobe.Output()
	pinData.Output()
	pinClock.Output()
	pinOE.Output()

	// What does this do??
	piOE.High()

	return rpio.Open()
}

func closeHardware() error {
	return rpio.Close()
}

func vendItem(location uint8) {
	setMotor(location)
}

// COPIED WHOLESALE FROM MOTORTEST - DOES IT WORK? WHO KNOWSSSS
func sendBit(state rpio.State) {
	pinClock.High()
	pinData.Write(state)
	pinClock.Low()
}

func setRegisters(r1 int) {
	pinStrobe.Low()
	for i := 0; i < 16; i++ {
		if r1 & 0x8000 {
			sendBit(rpio.High)
		} else {
			sendBit(rpio.Low)
		}
	}
	pinStrobe.High()
}

func setMotor(r1 int) {
	if r1 == 0 {
		setRegisters(0)
	} else {
		r1--
		setRegisters(1 << r1)
	}
}
