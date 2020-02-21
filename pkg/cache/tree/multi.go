package tree

type MultiTreeNode interface {
	Object() interface{}
	Children() []MultiTreeNode
	Insert(obj interface{}) MultiTreeNode
}

type multiTreeNode struct {

	// store object
	obj interface{}

	level uint64

	// the parent
	parent *multiTreeNode

	// the children
	children []*multiTreeNode
}

func (t *multiTreeNode) Object() interface{} {
	return t.obj
}

func (t *multiTreeNode) Children() []MultiTreeNode {

	var out []MultiTreeNode

	for _, e := range t.children {
		out = append(out, e)
	}
	return out
}

func (t *multiTreeNode) IsLeaf() bool {
	return len(t.children) == 0
}

func (t *multiTreeNode) IsRoot() bool {
	return t.parent == nil
}

// new a node for tree
func newMultiTreeNode(obj interface{}, parent *multiTreeNode) *multiTreeNode {
	return &multiTreeNode{
		obj:    obj,
		parent: parent,
		level:  parent.level + 1, // add level for the node
	}
}

func (t *multiTreeNode) Insert(obj interface{}) MultiTreeNode {

	out := newMultiTreeNode(obj, t)
	t.children = append(t.children, out)
	return out
}

// pre order traverse
// root -> left -> right
func (t *multiTreeNode) preOrder(action func(MultiTreeNode)) {

	// 1. visit the root element
	action(t)

	// 2. if it a leaf, just return
	if t.IsLeaf() {
		return
	}

	// 3. visit the child
	for _, c := range t.children {
		c.preOrder(action)
	}
}

// post order traverse
// left -> right -> root
func (t *multiTreeNode) postOrder(action func(MultiTreeNode)) {

	// 1. visit the leaf node
	if t.IsLeaf() {
		action(t)
		return
	}

	// 2. visit the child
	for _, c := range t.children {
		c.postOrder(action)
	}

	// 3. visit the root element
	action(t)
}

func (t *multiTreeNode) DRL(f func(MultiTreeNode)) {
	t.preOrder(f)
}

func (t *multiTreeNode) LRD(f func(MultiTreeNode)) {
	t.postOrder(f)
}
