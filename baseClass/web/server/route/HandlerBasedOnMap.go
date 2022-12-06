package route

import (
	"fmt"
	"net/http"
)

type HandlerBasedOnMap struct {
	Handlers map[string]func(c *Context)
}

func (h *HandlerBasedOnMap) Core(c *Context) {
	key := h.GetKey(c.R.Method, c.R.URL.Path)
	if handler, ok := h.Handlers[key]; !ok {
		handler(c)
	} else {
		c.W.WriteHeader(http.StatusNotFound)
		_, _ = c.W.Write([]byte("not any router match"))
		return
	}

}

func (h *HandlerBasedOnMap) Register(method string, pattern string, handlerFunc HandlerFunc) {
	key := h.GetKey(method, pattern)
	h.Handlers[key] = handlerFunc
}

func (h *HandlerBasedOnMap) GetKey(method, path string) string {
	return fmt.Sprintf("%s#%s", method, path)
}

func NewMapRoute() *HandlerBasedOnMap {
	return &HandlerBasedOnMap{}
}
