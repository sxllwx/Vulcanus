package types

// MultiTree
type MultiTree struct {

	// Tree Root
	Root *MultiTree
	// Parent Node
	Parent *MultiTree

	// actual object
	Item interface{}

	// tree dep
	Depth uint32

	// current node depth
	CurrentDepth uint32

	// children list
	// if children-list is nil, this a leaf
	ChildrenList []*MultiTree
}

func NewMultiTree(item interface{}) *MultiTree {

	n := &MultiTree{
		Item: item,
	}
	n.Root = n     // pointer to self
	n.Parent = nil // no parent
	return n
}

func (n *MultiTree) Insert(item interface{}) *MultiTree {

	cn := &MultiTree{
		Root:         n.Root,
		Parent:       n,
		Item:         item,
		CurrentDepth: n.CurrentDepth + 1,
	}

	if n.Root.Depth < cn.CurrentDepth {
		// update root depth
		n.Root.Depth = cn.CurrentDepth
	}

	n.ChildrenList = append(n.ChildrenList, cn)

	return cn
}

// deep-first traversal
func (n *MultiTree) deepTraversalChildrenList(f func(*MultiTree)) {

	f(n)

	for _, c := range n.ChildrenList {
		c.deepTraversalChildrenList(f)
	}

}

// breadth-first traversal
func (n *MultiTree) breadthTraversalChildrenList(f func(*MultiTree)) {

	// 1. visit child first
	for _, c := range n.ChildrenList {
		f(c)
	}

	// 2. recursive visit child's child
	for _, c := range n.ChildrenList {
		c.breadthTraversalChildrenList(f)
	}
}

func (n *MultiTree) BreadthFirstVisitChildrenList(f func(*MultiTree)) {
	n.breadthTraversalChildrenList(f)
}

func (n *MultiTree) DeepFirstVisitChildrenList(f func(*MultiTree)) {
	for _, c := range n.ChildrenList {
		c.deepTraversalChildrenList(f)
	}
}

func (n *MultiTree) VisitParent(f func(*MultiTree)) {

	if n.Parent == nil {
		// already is root
		return
	}

	f(n.Parent)
	n.Parent.VisitParent(f)
}
