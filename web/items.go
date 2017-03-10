package web

import (
	"net/http"

	"github.com/go-martini/martini"
	"github.com/lexicality/vending/backend"
	"github.com/martini-contrib/render"
	logging "github.com/op/go-logging"
)

func renderHome(
	req *http.Request,
	r render.Render,
	log *logging.Logger,
	stock *backend.Stock,
) {
	ctx := req.Context()
	items, err := stock.GetAll(ctx)
	if err != nil {
		log.Errorf("Unable to get homepage items: %s", err)
		r.HTML(500, "500", nil)
		return
	}
	r.HTML(200, "home", items)
}

func renderItem(
	req *http.Request,
	params martini.Params,
	r render.Render,
	log *logging.Logger,
	stock *backend.Stock,
) {
	ctx := req.Context()
	item, err := stock.GetItem(ctx, params["ID"])
	if err != nil {
		log.Errorf("Unable to get item details for %s: %s", params["ID"], err)
		r.HTML(500, "500", nil)
		return
	} else if item == nil {
		r.HTML(404, "404", nil)
		return
	}

	r.HTML(200, "item", item)
}
