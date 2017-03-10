package web

import (
	"github.com/go-martini/martini"
	"github.com/lexicality/vending/backend"
	"github.com/lexicality/vending/hardware"
	"github.com/lexicality/vending/vend"
	"github.com/martini-contrib/render"
	"github.com/op/go-logging"
)

func render404(r render.Render) {
	r.HTML(404, "404", nil)
}

func renderHome(r render.Render, log *logging.Logger, stock *backend.Stock) {
	items, err := stock.GetAll()
	if err != nil {
		log.Errorf("Unable to get homepage items: %s", err)
		r.HTML(500, "500", nil)
		return
	}
	r.HTML(200, "home", items)
}

func renderItem(params martini.Params, r render.Render, log *logging.Logger, stock *backend.Stock) {
	item, err := stock.GetItem(params["ID"])
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

// Server runs the web server (!)
func Server(addr, webRoot string, log *logging.Logger, stock *backend.Stock, hw hardware.Hardware) {
	m := martini.Classic()
	m.Use(render.Renderer(render.Options{
		Directory:  webRoot + "/tpl",
		Extensions: []string{".tmpl", ".html"},
		Layout:     "layout",
	}))
	m.Use(martini.Static(webRoot, martini.StaticOptions{
		Prefix:  "static",
		Exclude: "/static/tpl/",
	}))
	m.Map(stock)
	m.Map(log)
	m.Map(hw)

	m.Get("/", renderHome)
	m.Get("/items/:ID", renderItem)
	m.Get("/items/:ID/vend", renderVendItem)
	m.NotFound(render404)

	m.RunOnAddr(addr)
}
