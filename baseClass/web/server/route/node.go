package route

import "sort"

const (

	// 根节点，只有根用这个
	nodeTypeRoot = iota

	// *
	nodeTypeAny

	// 路径参数
	nodeTypeParam

	// 正则
	nodeTypeReg

	// 静态，即完全匹配
	nodeTypeStatic
)

func NewRootNode() *Node {
	return &Node{
		path:     "/",
		children: nil,
		handler:  nil,
	}
}

func NewStaticNode(path string) *Node {
	return &Node{
		path:     path,
		children: nil,
		handler:  nil,
		nodeType: nodeTypeAny,
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
		path:     path,
		children: nil,
		handler:  nil,
		nodeType: nodeTypeAny,
		matchFunc: func(p string) bool {
			return true
		},
	}
}

func NewAnyNode(path string) *Node {
	return &Node{
		path:     path,
		children: nil,
		handler:  nil,
		nodeType: nodeTypeAny,
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

type Node struct {
	path      string
	children  []*Node
	handler   HandlerFunc
	matchFunc // 判断path是否匹配
	nodeType  int
}

func (n *Node) CreateChildNode(root *Node, path string) *Node {
	if oldNode, ok := n.FindChildrenNode(root, path); ok {
		return oldNode
	}

	var newNode *Node
	if path == "*" {
		newNode = NewAnyNode(path)
	} else if path[:1] == ":" {
		newNode = NewAnyNode(path)
	} else {
		newNode = NewStaticNode(path)
	}

	root.children = append(root.children, newNode)
	return newNode
}

func (n *Node) FindChildrenNode(root *Node, path string) (*Node, bool) {
	for _, node := range root.children {
		if node.path == path {
			return node, true
		}
	}

	return nil, false

}

func (n *Node) MatchNode(root *Node, path string) (*Node, bool) {
	var result []*Node
	for _, node := range root.children {
		if node.matchFunc(path) {
			result = append(result, node)
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
