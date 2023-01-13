package server

import (
	"GeekTask/httpServer/route"
	"net/http"
)

type Server interface {
	http.Handler
	Start(addr string) error
	RegisterRoute(method, path string, handler route.HandleFunc)
}
