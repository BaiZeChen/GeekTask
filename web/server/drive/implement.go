package drive

import (
	"GeekTask/web/server/middleware"
	"GeekTask/web/server/route"
	"net/http"
)

type HttpServer struct {
	Name    string
	Handler route.Handler // 功能处理 （注册路由，发现路由）
	root    middleware.Middleware
}

func (h *HttpServer) Register(method, pattern string, handlerFunc route.HandlerFunc) {
	h.Handler.Register(method, pattern, handlerFunc)
}

func (h *HttpServer) Start(address string) error {
	return http.ListenAndServe("", h)
}

func (h *HttpServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	c := route.NewContext(writer, request)
	h.root(c)
}

func NewSdkHttpServer(name string, builders ...middleware.MiddlewareBuild) Server {

	// 改用我们的树
	handler := route.NewMapRoute()
	// 因为我们是一个链，所以我们把最后的业务逻辑处理，也作为一环
	var root middleware.Middleware = handler.Core
	// 从后往前把filter串起来
	for i := len(builders) - 1; i >= 0; i-- {
		b := builders[i]
		root = b(root)
	}
	res := &HttpServer{
		Name:    name,
		Handler: handler,
		root:    root,
	}
	return res
}
