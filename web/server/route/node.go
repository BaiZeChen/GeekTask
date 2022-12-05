package route

func NewRootNode() *Node {
	return &Node{
		path:     "/",
		children: nil,
		handler:  nil,
	}
}

type Node struct {
	path     string
	children []*Node
	handler  HandlerFunc
}

func (n *Node) CreateChildNode(root *Node, path string) *Node {
	if oldNode, ok := n.FindChildrenNode(root, path); ok {
		return oldNode
	}
	newNode := &Node{
		path:     path,
		children: nil,
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
