package server

import (
	"GeekTask/httpServer/route"
	"net/http"
)

type HttpServer struct {
	route *route.Router
	ms    []route.Middleware
}

func (h *HttpServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	ctx := &route.Context{
		Req:  request,
		Resp: writer,
	}
	root := h.handle
	msLen := len(h.ms)
	for i := msLen; i >= 0; i-- {
		root = h.ms[i](root)
	}
	root(ctx)
}

// 增加中间件
func (h *HttpServer) Use(ms ...route.Middleware) {
	h.ms = append(h.ms, ms...)
}

func (h *HttpServer) Start(addr string) error {
	return http.ListenAndServe(addr, h)
}

func (h *HttpServer) RegisterRoute(method, path string, handler route.HandleFunc) {
	err := h.route.Register(method, path, handler)
	if err != nil {
		panic(err)
	}
}

func (h *HttpServer) Post(path string, handler route.HandleFunc) {
	h.RegisterRoute(http.MethodPost, path, handler)
}

func (h *HttpServer) Get(path string, handler route.HandleFunc) {
	h.RegisterRoute(http.MethodGet, path, handler)
}

func (h *HttpServer) handle(ctx *route.Context) {
	routeHandle, ok := h.route.Find(ctx.Req.Method, ctx.Req.URL.Path, ctx)
	if !ok {
		ctx.Resp.WriteHeader(404)
		ctx.Resp.Write([]byte("Not Found"))
		return
	}
	routeHandle(ctx)
}
