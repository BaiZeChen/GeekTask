package route

type HandlerFunc func(c *Context)

// Routable 可路由的
type Routable interface {
	// Route 设定一个路由，命中该路由的会执行handlerFunc的代码
	Register(method string, pattern string, handlerFunc HandlerFunc)
}

type Handler interface {
	Core(c *Context)
	Routable
}
