package vending

import (
	"github.com/satori/go.uuid"
)

// Request represents a vending request from the server to the vending machine
type Request struct {
	// Location is the ID of the motorised arm thing it's in
	Location uint8 `json:"location"`
	// ID identifies the request to make reporting errors easier
	ID uuid.UUID `json:"id"`
}

func NewRequest(location uint8) *Request {
	return &Request{
		Location: location,
		ID:       uuid.NewV4(),
	}
}
