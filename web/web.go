package web

import (
	"context"

	"net/http"

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
func Server(
	ctx context.Context,
	addr string,
	webRoot string,
	log *logging.Logger,
	stock *backend.Stock,
	hw hardware.Hardware,
) error {
	doneC := ctx.Done()

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

	// Tell active HTTP requests to stop when we stop
	m.Use(func(req *http.Request, c martini.Context) {
		ctx := req.Context()
		newCtx, cancel := context.WithCancel(ctx)
		c.Map(req.WithContext(newCtx))
		go func() {
			select {
			case <-doneC:
				cancel()
			case <-ctx.Done():
				// exit
			}
		}()
	})
	m.Map(stock)
	m.Map(log)
	m.Map(hw)

	m.Get("/", renderHome)
	m.Get("/items/:ID", renderItem)
	m.Get("/items/:ID/vend", renderVendItem)
	m.NotFound(render404)

	server := &http.Server{
		Handler: m,
		Addr:    addr,
	}

	serverErrC := make(chan error)
	go func() {
		serverErrC <- server.ListenAndServe()
	}()

	select {
	case err := <-serverErrC:
		return err
	case <-doneC:
		// TODO: Timeouts?
		return server.Shutdown(context.TODO())
	}
}
