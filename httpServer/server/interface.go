package server

import (
	"GeekTask/baseClass/web/server/route"
	"net/http"
)

type HandleFunc func(ctx *route.Context)

type Server interface {
	http.Handler
	Start(addr string) error
	AddRoute(method, path string, handler HandleFunc)
}
