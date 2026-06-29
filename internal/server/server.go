package server

import (
	"net/http"
)

func New(handler http.Handler, addr string) *http.Server {
	return &http.Server{
		Addr:    addr,
		Handler: handler,
	}
}
