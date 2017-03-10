package web

import (
	"github.com/go-martini/martini"
	"github.com/lexicality/vending/backend"
	"github.com/lexicality/vending/hardware"
	"github.com/martini-contrib/render"
	"github.com/op/go-logging"
)

func render404(r render.Render) {
	r.HTML(404, "404", nil)
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
