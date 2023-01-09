package server

import (
	"GeekTask/httpServer/route"
	"net/http"
)

type HttpServer struct {
}

func (h *HttpServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	ctx := &route.Context{
		Req:  request,
		Resp: writer,
	}
	h.serve(ctx)

}

func (h *HttpServer) Start(addr string) error {
	return http.ListenAndServe(addr, h)
}

func (h *HttpServer) AddRoute(method, path string, handler HandleFunc) {
	panic("implement me")
}

func (h *HttpServer) serve(ctx *route.Context) {

}
