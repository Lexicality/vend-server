package backend

import (
	"context"
	"fmt"
	"sync"

	"github.com/lexicality/vending/vend"
)

// FakeStockItem is a data store for FakeStock
type FakeStockItem struct {
	StockItem
	sync.RWMutex
}

// GetStockItem returns a copy of the stock item
func (item *FakeStockItem) GetStockItem() *StockItem {
	item.RLock()
	ret := item.StockItem
	item.RUnlock()
	return &ret
}

// FakeStock is an in-memory Stock
type FakeStock struct {
	Items map[string]*FakeStockItem
}

func newStock() *FakeStock {
	return &FakeStock{
		Items: make(map[string]*FakeStockItem, 14),
	}
}

// GetFakeStock returns a Stock with a bunch of fake items in it
func GetFakeStock() *FakeStock {
	stock := newStock()

	stock.Items["c9b2d770-532f-41fc-acf7-e6e682bd82d5"] = &FakeStockItem{
		StockItem{
			ID:       "c9b2d770-532f-41fc-acf7-e6e682bd82d5",
			Name:     "Example Item #1",
			Quantity: 5,
			Price:    1000,
			Location: 0,
			Image:    "lorem-pixel-01.jpg",
		},
		sync.RWMutex{},
	}
	stock.Items["67fd2453-a2a9-4260-949d-c0f573b4b3ab"] = &FakeStockItem{
		StockItem{
			ID:       "67fd2453-a2a9-4260-949d-c0f573b4b3ab",
			Name:     "Example Item #2",
			Quantity: 50,
			Price:    333,
			Location: 1,
			Image:    "lorem-pixel-02.jpg",
		},
		sync.RWMutex{},
	}
	stock.Items["47db0c90-e9ff-4b2f-a646-ef259887a89b"] = &FakeStockItem{
		StockItem{
			ID:       "47db0c90-e9ff-4b2f-a646-ef259887a89b",
			Name:     "Example Item #3",
			Quantity: 200,
			Price:    11133,
			Location: 2,
			Image:    "lorem-pixel-03.jpg",
		},
		sync.RWMutex{},
	}
	stock.Items["aaac9bd8-ac3f-4348-af1c-c211396e3ff9"] = &FakeStockItem{
		StockItem{
			ID:       "aaac9bd8-ac3f-4348-af1c-c211396e3ff9",
			Name:     "Example Item #4",
			Quantity: 5,
			Price:    50,
			Location: 3,
			Image:    "lorem-pixel-04.jpg",
		},
		sync.RWMutex{},
	}
	stock.Items["552be5d3-ed40-4a2e-bb88-f0319b5f4af1"] = &FakeStockItem{
		StockItem{
			ID:       "552be5d3-ed40-4a2e-bb88-f0319b5f4af1",
			Name:     "Example Item #5",
			Quantity: 5,
			Price:    100,
			Location: 4,
			Image:    "lorem-pixel-05.jpg",
		},
		sync.RWMutex{},
	}
	stock.Items["ef8e1fc5-8a5d-49fb-baba-37377beead4a"] = &FakeStockItem{
		StockItem{
			ID:       "ef8e1fc5-8a5d-49fb-baba-37377beead4a",
			Name:     "Example Item #6",
			Quantity: 5,
			Price:    1200,
			Location: 5,
			Image:    "lorem-pixel-06.jpg",
		},
		sync.RWMutex{},
	}
	stock.Items["0840a34d-e024-4e47-a8db-471ac97c6aae"] = &FakeStockItem{
		StockItem{
			ID:       "0840a34d-e024-4e47-a8db-471ac97c6aae",
			Name:     "Example Item #7",
			Quantity: 5,
			Price:    123456,
			Location: 6,
			Image:    "lorem-pixel-07.jpg",
		},
		sync.RWMutex{},
	}
	stock.Items["193599ed-d1eb-411c-9ecc-e2343256609b"] = &FakeStockItem{
		StockItem{
			ID:       "193599ed-d1eb-411c-9ecc-e2343256609b",
			Name:     "Example Item #8",
			Quantity: 5,
			Price:    0,
			Location: 7,
			Image:    "lorem-pixel-08.jpg",
		},
		sync.RWMutex{},
	}
	stock.Items["c0e26c01-8f37-4a23-a7e9-caf322036ba9"] = &FakeStockItem{
		StockItem{
			ID:       "c0e26c01-8f37-4a23-a7e9-caf322036ba9",
			Name:     "Example Item #9",
			Quantity: 5,
			Price:    100,
			Location: 8,
			Image:    "lorem-pixel-09.jpg",
		},
		sync.RWMutex{},
	}
	stock.Items["b4d6ec03-268c-49d6-ac17-c107b3375014"] = &FakeStockItem{
		StockItem{
			ID:       "b4d6ec03-268c-49d6-ac17-c107b3375014",
			Name:     "Example Item #10",
			Quantity: 0,
			Price:    1000,
			Location: 9,
			Image:    "lorem-pixel-10.jpg",
		},
		sync.RWMutex{},
	}
	stock.Items["31d1fb41-6af1-4fff-ad2e-1138632a6dab"] = &FakeStockItem{
		StockItem{
			ID:       "31d1fb41-6af1-4fff-ad2e-1138632a6dab",
			Name:     "Example Item #11",
			Quantity: 5,
			Price:    1,
			Location: 10,
			Image:    "lorem-pixel-11.jpg",
		},
		sync.RWMutex{},
	}
	stock.Items["1bb541b6-b465-445e-8bb3-58688bf40e13"] = &FakeStockItem{
		StockItem{
			ID:       "1bb541b6-b465-445e-8bb3-58688bf40e13",
			Name:     "Example Item #12",
			Quantity: 8,
			Price:    100,
			Location: 11,
			Image:    "lorem-pixel-12.jpg",
			Broken:   true,
		},
		sync.RWMutex{},
	}
	return stock
}

func (stock *FakeStock) lookupItem(ID string) *FakeStockItem {
	return stock.Items[ID]
}

// GetAll returns all items currently in the vending machine (available or not)
func (stock *FakeStock) GetAll(ctx context.Context) (items []*StockItem, err error) {
	items = make([]*StockItem, len(stock.Items))
	i := 0
	for _, item := range stock.Items {
		items[i] = item.GetStockItem()
		i++
	}

	return
}

// GetItem returns information specific to a single item
func (stock *FakeStock) GetItem(ctx context.Context, ID string) (*StockItem, error) {
	item := stock.lookupItem(ID)
	if item != nil {
		return item.GetStockItem(), nil
	}
	return nil, nil
}

// ReserveItem indicates that you are queued to vend an item and it's unavailable for other vends
func (stock *FakeStock) ReserveItem(ctx context.Context, ID string) error {
	item := stock.lookupItem(ID)
	if item == nil {
		return ErrNotAnItem
	}

	item.RLock()
	if item.Quantity == 0 {
		return ErrItemEmpty
	}
	item.RUnlock()

	item.Lock()
	item.Reserved++
	item.Unlock()
	return nil
}

// UpdateItem updates the stock with the result of your recent vend attempt.
func (stock *FakeStock) UpdateItem(ctx context.Context, ID string, status vend.Result) error {
	item := stock.lookupItem(ID)
	if item == nil {
		return ErrNotAnItem
	}

	item.Lock()
	defer item.Unlock()

	switch status {
	case vend.ResultSuccess:
		item.Quantity--
		item.Reserved--
	case vend.ResultEmpty:
		item.Quantity = 0
		// TODO: Y'all need to get your money back somehow
		item.Reserved = 0
	case vend.ResultAborted:
		item.Reserved--
	case vend.ResultJammed:
	case vend.ResultHardwareFailure:
	case vend.ResultUnknownFailure:
		item.Broken = true
	default:
		return fmt.Errorf("unexpected result %v supplied to UpdateItem", status)
	}

	return nil
}
