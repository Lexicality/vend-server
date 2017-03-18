package web

import (
	"context"
	"fmt"
	"net/http"

	"github.com/martini-contrib/render"
	logging "github.com/op/go-logging"
	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/charge"

	"github.com/lexicality/vending/backend"
	"github.com/lexicality/vending/hardware"
	"github.com/lexicality/vending/vend"
)

func handleBuy(
	req *http.Request,
	r render.Render,
	log *logging.Logger,

	stock backend.Stock,
	hw hardware.Machine,
	txns backend.Transactions,
) {
	reqCtx := req.Context()
	globalCtx := reqCtx.Value(globalContextKey).(context.Context)

	// TODO: Parse & validate
	itemID := req.PostFormValue("item")
	user := req.PostFormValue("stripeEmail")
	token := req.FormValue("stripeToken")

	// Reserve Item
	err := stock.ReserveItem(reqCtx, itemID)
	if err == backend.ErrNotAnItem {
		r.Text(http.StatusBadRequest, "not item")
		return
	} else if err == backend.ErrItemBroken {
		// TODO: wat say + status code
		r.Text(http.StatusServiceUnavailable, "lp0 on fire")
		return
	} else if err == backend.ErrItemEmpty {
		// TODO: wat say + status code
		r.Text(http.StatusServiceUnavailable, "items empty soz")
		return
	} else if err != nil {
		log.Errorf("Unable to retrieve item %s: %s", itemID, err)
		r.HTML(500, "500", nil)
		return
	}
	abortReserve := func() {
		err := stock.UpdateItem(context.TODO(), itemID, vend.ResultAborted)
		if err != nil {
			log.Criticalf("Failed to abort item %s: %s", itemID, err)
		}
	}

	// Get the item (since reserve doesn't get it for us)
	item, err := stock.GetItem(reqCtx, itemID)
	if err != nil || item == nil {
		log.Errorf("Unable to retrieve item %s: %s", itemID, err)
		abortReserve()
		r.HTML(500, "500", nil)
		return
	}

	// Set up a transaction for this transaction
	txn, err := txns.New(globalCtx, item, user)
	if err != nil {
		log.Errorf("Unable to create transaction: %s", err)
		abortReserve()
		// Should probably mention how you've not been charged etc
		r.HTML(500, "500", nil)
		return
	}

	// Demand money from Stripe
	params := &stripe.ChargeParams{
		Amount:   item.Price,
		Currency: "gbp",
		Desc:     item.Name,
	}
	params.AddMeta("user", user)
	params.AddMeta("transaction", txn.ID.String())
	params.SetSource(token)

	// At this point reqCtx is inaproprate since this needs to continue even if the user closes the page
	// However we probably don't want to stop *now* since money is happening.
	charge, err := charge.New(params)

	if err != nil {
		// TODO: Actually check the error type etc etc
		txn.State = backend.TransactionFailed
		txn.Reason = "???"
		_ = txns.Update(context.TODO(), txn)
		log.Debugf("OMG ERROR %+v %s", err, err)
		abortReserve()
		r.Text(500, err.Error())
		return
	}

	log.Infof("%s has successfully paid %s for a %s via charge #%s", user, item.FormattedPrice(), item.Name, charge.ID)

	txn.ProviderID = charge.ID
	txn.State = backend.TransactionPaid
	// TODO: Context? Presumably we don't want to cancel this update given how important it is.
	txns.Update(context.TODO(), txn)

	if globalCtx.Err() != nil {
		// um
		log.Criticalf("Bailing out of incomplete vend %s due to context closing: %s", txn.ID, globalCtx.Err())
		r.Text(503, "Please contact the trustees")
		return
	}

	go func() {
		var err error
		res := hw.Vend(globalCtx, item.Location)
		if res == vend.ResultSuccess {
			txn.State = backend.TransactionComplete
			log.Noticef("Successful vend of %s", item.ID)
		} else {
			txn.State = backend.TransactionFailed
			txn.Reason = res.String()
			log.Criticalf("Unsuccessful vend of %s: %s", item.ID, res)
		}

		// TODO: Log errors
		err = stock.UpdateItem(context.TODO(), item.ID, res)
		if err != nil {
			log.Errorf("TODO UPDATE FAILED %s", err)
		}
		err = txns.Update(context.TODO(), txn)
		if err != nil {
			log.Errorf("TODO TRANSACTION FAILED %s", err)
		}
	}()

	r.Text(200, fmt.Sprintf("%+v", charge))
}
