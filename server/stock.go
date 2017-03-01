package main

import (
	"errors"

	"encoding/json"

	"github.com/joiggama/money"
	"github.com/lexicality/vending/shared/vending"
)

var (
	// ErrNotAnItem is what happens if you try to vend an item what not be there
	ErrNotAnItem = errors.New("Item does not exist")
	// ErrItemEmpty is if you vend something that's out
	ErrItemEmpty = errors.New("Item has no Quantity")
)

// StockItem represents an item in the vending machine
type StockItem struct {
	ID       string
	Name     string
	Quantity uint8
	Reserved uint8
	Image    string
	Price    int64
	Location uint8
}

var mFormatOptions = money.Options{"currency": "GBP"}

// FormattedPrice returns the price as a currency string
func (item *StockItem) FormattedPrice() string {
	return money.Format(float64(item.Price)/10, mFormatOptions)

}

// Stock is an interface to the vending machine (possibly could be named better)
type Stock struct {
	VendC <-chan string

	items map[string]*StockItem
	vendC chan string
}

func newStock() *Stock {
	stream := make(chan string)
	return &Stock{
		items: make(map[string]*StockItem, 14),
		vendC: stream,
		VendC: stream,
	}
}

// GetFakeStock returns a Stock with a bunch of fake items in it
func GetFakeStock() *Stock {
	stock := newStock()

	stock.items["c9b2d770-532f-41fc-acf7-e6e682bd82d5"] = &StockItem{
		ID:       "c9b2d770-532f-41fc-acf7-e6e682bd82d5",
		Name:     "Example Item #1",
		Quantity: 5,
		Image:    "",
		Price:    1000,
		Location: 0,
	}
	stock.items["67fd2453-a2a9-4260-949d-c0f573b4b3ab"] = &StockItem{
		ID:       "67fd2453-a2a9-4260-949d-c0f573b4b3ab",
		Name:     "Example Item #2",
		Quantity: 50,
		Price:    333,
		Location: 1,
	}
	stock.items["47db0c90-e9ff-4b2f-a646-ef259887a89b"] = &StockItem{
		ID:       "47db0c90-e9ff-4b2f-a646-ef259887a89b",
		Name:     "Example Item #3",
		Quantity: 200,
		Price:    11133,
		Location: 2,
	}
	stock.items["aaac9bd8-ac3f-4348-af1c-c211396e3ff9"] = &StockItem{
		ID:       "aaac9bd8-ac3f-4348-af1c-c211396e3ff9",
		Name:     "Example Item #4",
		Quantity: 5,
		Price:    50,
		Location: 3,
	}
	stock.items["552be5d3-ed40-4a2e-bb88-f0319b5f4af1"] = &StockItem{
		ID:       "552be5d3-ed40-4a2e-bb88-f0319b5f4af1",
		Name:     "Example Item #5",
		Quantity: 5,
		Price:    100,
		Location: 4,
	}
	stock.items["ef8e1fc5-8a5d-49fb-baba-37377beead4a"] = &StockItem{
		ID:       "ef8e1fc5-8a5d-49fb-baba-37377beead4a",
		Name:     "Example Item #6",
		Quantity: 5,
		Price:    1200,
		Location: 5,
	}
	stock.items["0840a34d-e024-4e47-a8db-471ac97c6aae"] = &StockItem{
		ID:       "0840a34d-e024-4e47-a8db-471ac97c6aae",
		Name:     "Example Item #7",
		Quantity: 5,
		Price:    123456,
		Location: 6,
	}
	stock.items["193599ed-d1eb-411c-9ecc-e2343256609b"] = &StockItem{
		ID:       "193599ed-d1eb-411c-9ecc-e2343256609b",
		Name:     "Example Item #8",
		Quantity: 5,
		Price:    0,
		Location: 7,
	}
	stock.items["c0e26c01-8f37-4a23-a7e9-caf322036ba9"] = &StockItem{
		ID:       "c0e26c01-8f37-4a23-a7e9-caf322036ba9",
		Name:     "Example Item #9",
		Quantity: 5,
		Price:    -100,
		Location: 8,
	}
	stock.items["b4d6ec03-268c-49d6-ac17-c107b3375014"] = &StockItem{
		ID:       "b4d6ec03-268c-49d6-ac17-c107b3375014",
		Name:     "Example Item #10",
		Quantity: 0,
		Price:    1000,
		Location: 9,
	}
	return stock
}

// GetAll returns all items currently in the vending machine (available or not)
func (stock *Stock) GetAll() (items []*StockItem, err error) {
	items = make([]*StockItem, len(stock.items))
	i := 0
	for _, item := range stock.items {
		items[i] = item
		i++
	}

	return
}

// GetItem returns information specific to a single item
func (stock *Stock) GetItem(ID string) (item *StockItem, err error) {
	return stock.items[ID], nil
}

// VendItem vends an item
func (stock *Stock) VendItem(ID string) (result vending.Result, err error) {
	item, err := stock.GetItem(ID)
	if err != nil {
		return vending.NoResult, err
	} else if item == nil {
		return vending.NoResult, ErrNotAnItem
	} else if item.Quantity == 0 || item.Quantity <= item.Reserved {
		return vending.ResultEmpty, nil
	}

	log.Infof("Vending %s which is at %d", item.Name, item.Location)

	req, err := json.Marshal(vending.NewRequest(item.Location))
	if err != nil {
		return
	}
	stock.vendC <- string(req)
	// TODO: Verification, reliability, etc

	return vending.ResultSuccess, nil
}
