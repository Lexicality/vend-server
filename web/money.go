package web

import (
	"net/http"

	"github.com/martini-contrib/render"
	logging "github.com/op/go-logging"
	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/charge"

	"fmt"

	"github.com/lexicality/vending/backend"
)

func handleBuy(r render.Render, stock *backend.Stock, req *http.Request, log *logging.Logger) {
	ctx := req.Context()

	// TODO: Parse & validate
	itemID := req.PostFormValue("item")
	user := req.PostFormValue("stripeEmail")
	token := req.FormValue("stripeToken")

	item, err := stock.GetItem(ctx, itemID)
	if err != nil {
		log.Errorf("Unable to retrieve item %s: %s", itemID, err)
		r.HTML(500, "500", nil)
		return
	} else if item == nil {
		log.Warningf("Got payment request from %s for non-existant item %s", user, itemID)
		r.Text(400, "??? Missing item")
		return
	}

	params := &stripe.ChargeParams{
		Amount:   item.Price,
		Currency: "gbp",
		Desc:     item.Name,
	}
	params.AddMeta("user", user)
	params.SetSource(token)

	charge, err := charge.New(params)

	if err != nil {
		log.Debugf("OMG ERROR %+v %s", err, err)
		r.Text(500, err.Error())
		return
	}

	r.Text(200, fmt.Sprintf("%+v", charge))
}
