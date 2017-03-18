package backend

import (
	"context"
	"fmt"
	"sync"
	"time"

	uuid "github.com/satori/go.uuid"
)

// FakeTransaction is the memory backing store for FakeTransactions
type FakeTransaction struct {
	Transaction
	sync.RWMutex
}

// GetTransaction returns a copy of the transaction
func (txn *FakeTransaction) GetTransaction() *Transaction {
	txn.RLock()
	ret := txn.Transaction
	txn.RUnlock()
	return &ret
}

// FakeTransactions handles syncing transactions with the database etc
type FakeTransactions struct {
	txns map[uuid.UUID]*FakeTransaction
}

// NewFakeTransactions does the needful
func NewFakeTransactions(ctx context.Context) *FakeTransactions {
	return &FakeTransactions{
		txns: make(map[uuid.UUID]*FakeTransaction),
	}
}

func (txmgr *FakeTransactions) lookup(ID string) (*FakeTransaction, error) {
	uid, err := uuid.FromString(ID)
	if err != nil {
		return nil, err
	}

	return txmgr.txns[uid], nil
}

// Get retrieves an existing transaction
func (txmgr *FakeTransactions) Get(ctx context.Context, ID string) (*Transaction, error) {
	txn, err := txmgr.lookup(ID)

	if err != nil {
		return nil, err
	} else if txn == nil {
		return nil, nil
	}

	return txn.GetTransaction(), nil
}

// Update stores the new transaction state in the database
func (txmgr *FakeTransactions) Update(ctx context.Context, txn *Transaction) error {
	ftxn, ok := txmgr.txns[txn.ID]
	if !ok {
		return fmt.Errorf("invalid transaction %s", txn.ID)
	}

	ftxn.Lock()
	ftxn.ProviderID = txn.ProviderID
	ftxn.State = txn.State
	ftxn.Reason = txn.Reason
	ftxn.Unlock()

	return nil
}

// New sets up a new transaction starting now
func (txmgr *FakeTransactions) New(ctx context.Context, item *StockItem, user string) (*Transaction, error) {
	txn := &FakeTransaction{
		Transaction: Transaction{
			ID:     uuid.NewV4(),
			Item:   item.ID,
			Amount: item.Price,
			User:   user,
			State:  TransactionPending,
			Date:   time.Now(),
		},
	}

	txmgr.txns[txn.ID] = txn

	return txn.GetTransaction(), nil
}
