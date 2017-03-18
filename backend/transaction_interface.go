package backend

import (
	"context"
	"time"

	uuid "github.com/satori/go.uuid"
)

// TransactionState describes the state of a transaction
type TransactionState uint8

const (
	// TransactionPending is when we have a token from Stripe but it's not been charged
	TransactionPending TransactionState = iota + 1

	// TransactionPaid is when stripe has confirmed the charge but we've not yet vended
	TransactionPaid

	// TransactionComplete is when we've successfully vended the item
	TransactionComplete

	// TransactionRejected is when Stripe declines us
	TransactionRejected
	// TransactionFailed is when something went wrong and we didn't vend
	TransactionFailed
)

// Transaction represents a single transaction in the system
type Transaction struct {
	ID         uuid.UUID
	Item       string
	Amount     uint64
	ProviderID string
	User       string // Name + email of person doing this
	State      TransactionState
	Date       time.Time
	Reason     string // If something went wrong - what
}

// Transactions looks after transactions
type Transactions interface {
	// New creates a transaction from the passed info and returns it
	New(ctx context.Context, item *StockItem, user string) (*Transaction, error)
	// Get returns an existing transaction or nil if it doesn't exist
	Get(ctx context.Context, ID string) (*Transaction, error)
	// Update updates an existing transaction to the new state.
	// Note: It will only update the fields ProviderID, State and Reason.
	Update(ctx context.Context, txn *Transaction) error
}
