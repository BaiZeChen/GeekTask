package server

import (
	"GeekTask/httpServer/middleware"
	"GeekTask/httpServer/route"
	"errors"
	"net/http"
)

type HttpServer struct {
	route         *route.Router
	globalMiddles []route.Middleware
}

func (h *HttpServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	ctx := &route.Context{
		Req:  request,
		Resp: writer,
	}
	root := h.handle
	middles := h.getMiddles(ctx)
	if len(middles) > 0 {
		h.globalMiddles = append(h.globalMiddles, middles...)
	}
	msLen := len(h.globalMiddles)
	for i := msLen; i >= 0; i-- {
		root = h.globalMiddles[i](root)
	}

	// 系统中间件
	root = middleware.NewRecover(500).Build()(root)
	root = middleware.NewFlashResp().Build()(root)
	root(ctx)

}

// 增加全局中间件
func (h *HttpServer) Use(ms ...route.Middleware) {
	h.globalMiddles = append(h.globalMiddles, ms...)
}

func (h *HttpServer) Start(addr string) error {
	return http.ListenAndServe(addr, h)
}

func (h *HttpServer) RegisterRoute(method, path string, handler route.HandleFunc, middlewares ...route.Middleware) {
	err := h.route.Register(method, path, handler)
	if err != nil {
		panic(err)
	}
	if len(middlewares) > 0 {
		err := h.route.RegisterMiddleWare(path, middlewares...)
		if err != nil {
			panic(err)
		}
	}
}

func (h *HttpServer) registerMiddle(path string, middlewares ...route.Middleware) {
	if len(middlewares) == 0 {
		panic(errors.New("中间件数量不能为空"))
	}
	err := h.route.RegisterMiddleWare(path, middlewares...)
	if err != nil {
		panic(err)
	}
}

func (h *HttpServer) Post(path string, handler route.HandleFunc, middlewares ...route.Middleware) {
	h.RegisterRoute(http.MethodPost, path, handler, middlewares...)
}

func (h *HttpServer) Get(path string, handler route.HandleFunc, middlewares ...route.Middleware) {
	h.RegisterRoute(http.MethodGet, path, handler, middlewares...)
}

func (h *HttpServer) handle(ctx *route.Context) {
	routeHandle, ok := h.route.FindHandle(ctx.Req.Method, ctx.Req.URL.Path, ctx)
	if !ok {
		ctx.RespStatusCode = http.StatusNotFound
		return
	}
	routeHandle(ctx)
}

func (h *HttpServer) getMiddles(ctx *route.Context) []route.Middleware {
	return h.route.FindMiddle(ctx.Req.URL.Path)
}
func (h *HttpServer) Group(prefix string) *Group {
	return &Group{
		prefix: prefix,
		s:      h,
	}
}

type Group struct {
	prefix string
	s      *HttpServer
}

func (g *Group) Use(path string, middlewares ...route.Middleware) *Group {
	path = g.prefix + path
	g.s.registerMiddle(path, middlewares...)
	return g
}

func (g *Group) Post(path string, handler route.HandleFunc, middlewares ...route.Middleware) {
	path = g.prefix + path
	g.s.Post(path, handler, middlewares...)
}

func (g *Group) Get(path string, handler route.HandleFunc, middlewares ...route.Middleware) {
	path = g.prefix + path
	g.s.Get(path, handler, middlewares...)
}
