package web

import (
	"github.com/go-martini/martini"
	"github.com/lexicality/vending/backend"
	"github.com/lexicality/vending/hardware"
	"github.com/lexicality/vending/vend"
	"github.com/martini-contrib/render"
	logging "github.com/op/go-logging"
)

type vendRenderdata struct {
	Item    *backend.StockItem
	Result  vend.Result
	Results map[string]vend.Result
}

func renderVendItem(params martini.Params, r render.Render, log *logging.Logger, stock *backend.Stock, hw hardware.Hardware) {
	item, err := stock.GetItem(params["ID"])
	if err != nil {
		log.Errorf("Unable to get item details for %s: %s", params["ID"], err)
		r.HTML(500, "500", nil)
		return
	} else if item == nil {
		r.HTML(404, "404", nil)
		return
	}

	var result vend.Result
	if item.CanVend() {
		// TODO: This ignores hardware availability etc
		result = hw.Vend(item.Location)
	} else {
		result = vend.ResultEmpty
	}

	r.HTML(200, "vend", &vendRenderdata{
		Item:    item,
		Result:  result,
		Results: vend.AllResults,
	})
}
