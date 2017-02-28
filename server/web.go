package main

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
)

func webServer(addr, webRoot string) {
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

	m.Get("/", func(r render.Render) {
		r.HTML(200, "home", nil)
	})

	m.Get("/test", func(r render.Render) {
		r.HTML(200, "test", nil)
	})

	m.Get("/ws", wsHandler)
	m.RunOnAddr(":8080")
}
