package main

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
)

func render404(r render.Render) {
	r.HTML(404, "404", nil)
}

func renderHome(r render.Render, stock *Stock) {
	items, err := stock.GetAll()
	if err != nil {
		log.Errorf("Unable to get homepage items: %s", err)
		r.HTML(500, "500", nil)
		return
	}
	r.HTML(200, "home", items)
}

func renderItem(params martini.Params, r render.Render, stock *Stock) {
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

func webServer(addr, webRoot string, stock *Stock) {
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

	m.Get("/test", func(r render.Render) {
		r.HTML(200, "test", nil)
	})

	m.Get("/", renderHome)
	m.Get("/items/:ID", renderItem)
	m.Get("/ws", wsHandler)
	m.NotFound(render404)

	m.RunOnAddr(":8080")
}
