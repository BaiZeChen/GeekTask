package route

const (
	NODE_ROOT = iota
	NODE_ANY
	NODE_PARA
	NODE_REG
	NODE_STATIC
)

func NewRootNode() *Node {
	return &Node{
		path:     "/",
		nodeType: NODE_ROOT,
	}
}

func NewStaticNode(path string) *Node {
	return &Node{
		path:     path,
		nodeType: NODE_STATIC,
		matchFunc: func(p string) bool {
			if p == path {
				return true
			} else {
				return false
			}
		},
	}
}

func NewParamNode(path string) *Node {
	return &Node{
		path: path,
		matchFunc: func(p string) bool {
			return true
		},
	}
}

func NewAnyNode(path string) *Node {
	return &Node{
		path:     path,
		nodeType: NODE_ANY,
		matchFunc: func(p string) bool {
			if path == "*" {
				return true
			} else {
				return false
			}
		},
	}
}

type matchFunc func(path string) bool

type HandleFunc func(ctx *Context)

type Node struct {
	nodeType int
	path     string
	children []*Node
	matchFunc
	handler HandleFunc
}
