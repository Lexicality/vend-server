package vending

import (
	"encoding/json"

	"github.com/satori/go.uuid"
)

// Message is a generic message that explains what it is
type Message struct {
	Type    string          `json:"type"`
	Message json.RawMessage `json:"message"`
}

// Request represents a vending request from the server to the vending machine
type Request struct {
	// Location is the ID of the motorised arm thing it's in
	Location uint8 `json:"location"`
	// ID identifies the request to make reporting errors easier
	ID uuid.UUID `json:"id"`
}

// NewRequestMessage creates a JSON encoded request with a random UUID
func NewRequestMessage(location uint8) ([]byte, error) {
	req, err := json.Marshal(&Request{
		Location: location,
		ID:       uuid.NewV4(),
	})
	if err != nil {
		return nil, err
	}

	return json.Marshal(&Message{
		Type:    "Request",
		Message: req,
	})
}
