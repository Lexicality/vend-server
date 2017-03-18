package backend

import (
	"context"
	"errors"

	"github.com/joiggama/money"
	"github.com/lexicality/vending/vend"
)

var (
	// ErrNotAnItem is what happens if you try to vend an item what not be there
	ErrNotAnItem = errors.New("invalid item ID")
	// ErrItemEmpty is if you vend something that's out
	ErrItemEmpty = errors.New("no stock available")
	// ErrItemBroken in if you vend something that's
	ErrItemBroken = errors.New("item is jammed")
)

// StockItem represents the current state of an item in the vending machine.
// It should not be retained as any changes to the stock (eg vends) will not be propagated
type StockItem struct {
	ID       string
	Name     string
	Quantity uint8
	Reserved uint8
	Image    string
	Price    uint64
	Location uint8
	Broken   bool
}

// CanVend checks stock availability
func (item *StockItem) CanVend() bool {
	q := item.Quantity
	r := item.Reserved
	// We can only vend if there are unreserved items available
	return q > 0 && q > r && !item.Broken
}

var mFormatOptions = money.Options{"currency": "GBP"}

// FormattedPrice returns the price as a currency string
func (item *StockItem) FormattedPrice() string {
	return money.Format(float64(item.Price)/100, mFormatOptions)
}

// Stock is a storage container for items in the vending machine
type Stock interface {
	// GetAll returns all items currently in the vending machine (available or not)
	GetAll(context.Context) ([]*StockItem, error)
	// GetItem returns information specific to a single item
	GetItem(ctx context.Context, ID string) (*StockItem, error)
	// ReserveItem indicates that you are queued to vend an item and it's unavailable for other vends
	ReserveItem(ctx context.Context, ID string) error
	// UpdateItem updates the stock with the result of your recent vend attempt.
	UpdateItem(ctx context.Context, ID string, status vend.Result) error
}
