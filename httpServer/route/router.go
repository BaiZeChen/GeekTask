package route

import (
	"errors"
	"sort"
	"strings"
)

type Router struct {
	tree   map[string]*Node
	msTree *Node
}

func (r *Router) Register(method, path string, handler HandleFunc) error {
	if _, ok := r.tree[method]; !ok {
		r.tree[method] = NewRootNode()
	}
	front := r.tree[method]

	pathSlice := strings.Split(strings.Trim(path, "/"), "/")
	node, err := r.createChildNode(front, pathSlice)
	if err != nil {
		return err
	}
	node.handler = handler

	return nil
}

func (r *Router) RegisterMiddleWare(path string, middlewares ...Middleware) error {
	middlewaresLen := len(middlewares)
	if middlewaresLen == 0 {
		return nil
	}
	if r.msTree == nil {
		r.msTree = NewRootNode()
	}
	front := r.msTree

	pathSlice := strings.Split(strings.Trim(path, "/"), "/")
	node, err := r.createChildNode(front, pathSlice)
	if err != nil {
		return err
	}
	node.middlewares = middlewares

	return nil
}

func (r *Router) FindHandle(method, pattern string, ctx *Context) (HandleFunc, bool) {
	var (
		is_find bool
		handle  HandleFunc
	)
	defer func() {
		if !is_find {
			ctx.Params = make(map[string]string)
		}
	}()

	var root *Node
	if _, ok := r.tree[method]; !ok {
		return nil, is_find
	}

	root = r.tree[method]
	patternSlice := strings.Split(strings.Trim(pattern, "/"), "/")
	for _, s := range patternSlice {
		node, find := r.matchNode(root, s)
		if !find {
			return nil, is_find
		}
		if node.nodeType == NODE_PARA {
			ctx.Params[node.path[1:]] = s
		}
		handle = node.handler
	}

	if handle == nil {
		return nil, is_find
	} else {
		is_find = true
		return handle, is_find
	}

}

func (r *Router) FindMiddle(pattern string) []Middleware {
	var (
		ms []Middleware
	)

	var root *Node
	if r.msTree == nil {
		return nil
	}

	root = r.msTree
	patternSlice := strings.Split(strings.Trim(pattern, "/"), "/")
	for _, s := range patternSlice {
		node, find := r.matchNode(root, s)
		if !find {
			return nil
		}
		ms = node.middlewares
	}
	return ms
}

func (r *Router) createChildNode(root *Node, paths []string) (*Node, error) {
	cue := root
	for _, path := range paths {
		var (
			node *Node
			ok   bool
			err  error
		)
		node, ok, err = r.findChildrenNode(cue, path)
		if err != nil {
			return nil, err
		}
		if !ok {
			if path == "*" {
				node = NewAnyNode(path)
				cue.children = append(cue.children, node)
			} else if path[:1] == ":" {
				node = NewParamNode(path)
				cue.children = append(cue.children, node)
			} else {
				node = NewStaticNode(path)
				cue.children = append(cue.children, node)
			}
		}
		cue = node
	}

	return cue, nil
}

func (r *Router) findChildrenNode(root *Node, path string) (*Node, bool, error) {
	for _, child := range root.children {
		// 针对参数路由特殊匹配处理
		if path[:1] == ":" {
			if child.nodeType == NODE_PARA {
				return nil, false, errors.New("参数路由不允许重复")
			}
		} else if (path[:1] == ":" && child.nodeType == NODE_ANY) || (path == "*" && child.nodeType == NODE_PARA) {
			return nil, false, errors.New("参数路由不允许重复")
		} else {
			if child.path == path {
				return child, true, nil
			}
		}
	}
	return nil, false, nil
}

func (r *Router) matchNode(root *Node, path string) (*Node, bool) {
	var result []*Node
	for _, children := range root.children {
		if children.matchFunc(path) {
			result = append(result, children)
		}
	}
	if len(result) == 0 {
		return nil, false
	}
	// 最后选出匹配度最高的
	sort.Slice(result, func(i, j int) bool {
		return result[i].nodeType < result[j].nodeType
	})

	return result[len(result)-1], true
}
