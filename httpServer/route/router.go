package route

import (
	"errors"
	"sort"
	"strings"
)

type Router struct {
	tree map[string]*Node
}

func (r *Router) Register(method, path string, handler HandleFunc) error {
	if _, ok := r.tree[method]; !ok {
		r.tree[method] = NewRootNode()
	}
	front := r.tree[method]

	pathSlice := strings.Split(strings.Trim(path, "/"), "/")
	err := r.createChildNode(front, pathSlice, handler)
	if err != nil {
		return err
	}

	return nil
}

func (r *Router) Find(method, pattern string, ctx *Context) (HandleFunc, bool) {
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

func (r *Router) createChildNode(root *Node, paths []string, handler HandleFunc) error {
	cue := root
	for _, path := range paths {
		if node, ok := r.findChildrenNode(cue, path); ok {
			cue = node
			// /user/:id  和 /user/:name不能同时出现
			if path[:1] == ":" {
				return errors.New("参数路由只能注册一次")
			}
		} else {
			var node *Node
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
			cue = node
		}
	}

	cue.handler = handler
	return nil
}

func (r *Router) findChildrenNode(root *Node, path string) (*Node, bool) {
	for _, child := range root.children {
		// 针对参数路由特殊匹配处理
		if path[:1] == ":" {
			if child.nodeType == NODE_PARA {
				return child, true
			}
		} else {
			if child.path == path {
				return child, true
			}
		}
	}
	return nil, false
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
