package web

import (
	"net/http"

	"github.com/gorilla/handlers"
)

var handler http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(400)
}

// ServeCanonical makes sure everything goes where it should
func ServeCanonical(addr, host string) error {
	return http.ListenAndServe(addr, handlers.CanonicalHost(host, 301)(handler))
}
