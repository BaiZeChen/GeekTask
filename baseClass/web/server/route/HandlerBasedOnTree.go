package route

import (
	"errors"
	"net/http"
	"strings"
)

type HandlerBasedOnTree struct {
	root map[string]*Node
}

func (h *HandlerBasedOnTree) Core(c *Context) {
	node, err := h.match(h.root[c.R.Method], c)
	if err != nil {
		c.W.WriteHeader(http.StatusNotFound)
		_, _ = c.W.Write([]byte("not any router match"))
		return
	}
	node.handler(c)
}

func (h *HandlerBasedOnTree) Register(method string, pattern string, handlerFunc HandlerFunc) {
	h.add(method, pattern, handlerFunc)
}

func (h *HandlerBasedOnTree) add(method string, pattern string, handlerFunc HandlerFunc) {
	if _, ok := h.root[method]; !ok {
		h.root[method] = NewRootNode()
	}
	front := h.root[method]

	// 将pattern按照URL的分隔符切割
	// 例如，/user/friends 将变成 [user, friends]
	// 将前后的/去掉，统一格式
	pattern = strings.Trim(pattern, "/")
	paths := strings.Split(pattern, "/")
	for _, path := range paths {
		front = front.CreateChildNode(front, path)
	}
	front.handler = handlerFunc

}

func (h *HandlerBasedOnTree) match(root *Node, c *Context) (*Node, error) {

	pattern := c.R.URL.Path
	// 去除头尾可能有的/，然后按照/切割成段
	front := root
	paths := strings.Split(strings.Trim(pattern, "/"), "/")
	var ok = true
	for _, path := range paths {
		front, ok = root.FindChildrenNode(front, path)
		if !ok {
			return nil, errors.New("没有找到对应的路由")
		}
		if front.nodeType == nodeTypeParam {
			c.PathParams[front.path[1:]] = path
		}
	}
	if front.handler == nil {
		return nil, errors.New("没有找到对应的路由")
	}

	return front, nil

}

func NewTreeRoute() *HandlerBasedOnTree {
	return &HandlerBasedOnTree{}
}
