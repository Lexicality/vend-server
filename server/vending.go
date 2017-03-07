package main

import (
	"crypto/x509"

	"errors"

	"github.com/lexicality/vending/shared/vending"
)

var (
	// ErrVendTimedOut is triggered when a request is sent but no response is recieved
	ErrVendTimedOut = errors.New("Timed out waiting for confirmation")
)

// VendingMachine represents the physical machine in the space
type VendingMachine struct {
	Name       string // Name on certificate
	Connected  bool   // Is the websocket active
	Functional bool   // Has the machine reported a hardware failure

	// Transactions
	waitingOp *vending.Request

	// Connectivity
	messageTo    chan<- *vending.SendMessage
	messageFrom  <-chan *vending.RecvMessage
	disconnected <-chan bool
}

// ValidateCert makes sure the connecting client is actually the machine
func (vm *VendingMachine) ValidateCert(cert *x509.Certificate) bool {
	return cert.Subject.CommonName == vm.Name
}

// Vend requests the machine do it's thing, blocking until a result is recieved
func (vm *VendingMachine) Vend(location uint8) (vending.Result, error) {
	return vending.NoResult, ErrVendTimedOut
}
