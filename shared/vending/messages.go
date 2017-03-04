package vending

import (
	"encoding/json"

	"github.com/satori/go.uuid"
)

// MessageProtocol is the Websocket protocol that describes using these Message structs to communicate
const MessageProtocol = "vend-json"

// SendMessage reperesents a Message being sent
// Message is an arbitrary JSON object, while Type describes how to handle it
type SendMessage struct {
	Type    string      `json:"type"`
	Message interface{} `json:"message"`
}

// RecvMessage represents a Message being recieved
// Type describes how to decode the JSON stored in Message
type RecvMessage struct {
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

// NewMessageID returns a unique identifier
func NewMessageID() uuid.UUID {
	return uuid.NewV4()
}

// CompareMessageIDs checks if two message IDs are equal
func CompareMessageIDs(one, two uuid.UUID) bool {
	return uuid.Equal(one, two)
}
