package web

import (
	"context"
	"crypto/tls"

	"net/http"

	"github.com/go-martini/martini"
	"github.com/lexicality/vending/backend"
	"github.com/lexicality/vending/hardware"
	"github.com/martini-contrib/render"
	"github.com/op/go-logging"
)

// Server represents the web server
type Server struct {
	Addr        string
	ServerName  string
	WebRoot     string
	TLSCertFile string
	TLSKeyFile  string
}

func render404(r render.Render) {
	r.HTML(404, "404", nil)
}

// ServeHTTP runs the web server (!)
func (srv *Server) ServeHTTP(
	ctx context.Context,
	log *logging.Logger,
	stock *backend.Stock,
	hw hardware.Hardware,
) error {
	// Set up tls
	var tlsCFG *tls.Config
	if srv.TLSCertFile != "" && srv.TLSKeyFile != "" {
		cert, err := tls.LoadX509KeyPair(srv.TLSCertFile, srv.TLSKeyFile)
		if err != nil {
			return err
		}

		tlsCFG = &tls.Config{
			Certificates: []tls.Certificate{cert},
			ServerName:   srv.ServerName,
		}
	}

	doneC := ctx.Done()

	// Set up Martini
	m := martini.Classic()
	m.Use(render.Renderer(render.Options{
		Directory:  srv.WebRoot + "/tpl",
		Extensions: []string{".tmpl", ".html"},
		Layout:     "layout",
	}))
	m.Use(martini.Static(srv.WebRoot, martini.StaticOptions{
		Prefix:      "static",
		Exclude:     "/static/tpl/",
		SkipLogging: true,
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

	// Set up the route handlers
	m.Get("/", renderHome)
	m.Get("/items/:ID", renderItem)
	m.Get("/items/:ID/vend", renderVendItem)
	m.Post("/vend", handleBuy)
	m.NotFound(render404)

	// Run the actual server
	server := &http.Server{
		Handler:   m,
		Addr:      srv.Addr,
		TLSConfig: tlsCFG,
	}

	serverErrC := make(chan error)
	go func() {
		var err error
		if tlsCFG != nil {
			err = server.ListenAndServeTLS("", "")
		} else {
			err = server.ListenAndServe()
		}
		serverErrC <- err
	}()

	// Listen for shutdowns
	select {
	case err := <-serverErrC:
		return err
	case <-doneC:
		// TODO: Timeouts?
		return server.Shutdown(context.TODO())
	}
}
