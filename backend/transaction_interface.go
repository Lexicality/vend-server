package backend

import (
	"context"
	"encoding/json"
	"strings"
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

func (st TransactionState) String() string {
	switch st {
	case TransactionPending:
		return "Pending"
	case TransactionPaid:
		return "Paid"
	case TransactionComplete:
		return "Complete"
	case TransactionRejected:
		return "Rejected"
	case TransactionFailed:
		return "Failed"
	default:
		return "Invalid"
	}
}

// MarshalJSON returns the value of String but lowercase
func (st TransactionState) MarshalJSON() ([]byte, error) {
	v := strings.ToLower(st.String())
	return json.Marshal(v)
}

// Transaction represents a single transaction in the system
// It has a very minimal JSON footprint to improve transfers
type Transaction struct {
	ID         uuid.UUID        `json:"id"`
	Item       string           `json:"-"`
	Amount     uint64           `json:"-"`
	ProviderID string           `json:"-"`
	User       string           `json:"-"` // Name + email of person doing this
	State      TransactionState `json:"state"`
	Date       time.Time        `json:"-"`
	Reason     string           `json:"reason,omitempty"` // If something went wrong - what
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
