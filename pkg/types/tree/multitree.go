package tree

import (
	"encoding/json"
)

// MultiTree
// support json marshal && unmarshal the tree to store
type MultiTree struct {

	// Tree Root
	Root *MultiTree `json:"-"`
	// Parent Node
	Parent *MultiTree `json:"-"`

	// actual object
	Item interface{} `json:"item"`

	// tree dep
	// shared all node in the tree
	Depth *uint32 `json:"depth"`

	// current node depth
	CurrentDepth uint32 `json:"current_depth"`

	// children list
	// if children-list is nil, this a leaf
	ChildrenList []*MultiTree `json:"children_list"`
}

// NewMultiTree
func NewMultiTree(item interface{}) *MultiTree {

	var dep uint32

	n := &MultiTree{
		Item:  item,
		Depth: &dep,
	}
	n.Root = n     // pointer to self
	n.Parent = nil // no parent
	return n
}

// Insert
// insert a element for spec tree
func (n *MultiTree) Insert(item interface{}) *MultiTree {

	cn := &MultiTree{
		Root:         n.Root,
		Parent:       n,
		Item:         item,
		CurrentDepth: n.CurrentDepth + 1,
		Depth:        n.Depth,
	}

	if *n.Root.Depth < cn.CurrentDepth {
		// update root depth
		*n.Root.Depth = cn.CurrentDepth
	}

	n.ChildrenList = append(n.ChildrenList, cn)

	return cn
}

// deep-first traverse
func (n *MultiTree) deepTraverseChildrenList(f func(*MultiTree) bool) {

	if f(n) {
		return
	}

	for _, c := range n.ChildrenList {
		c.deepTraverseChildrenList(f)
	}

}

// breadth-first traverse
func (n *MultiTree) breadthTraverseChildrenList(f func(*MultiTree) bool) {

	// 1. visit child first
	for _, c := range n.ChildrenList {
		if f(c) {
			return
		}
	}

	// 2. recursive visit child's child
	for _, c := range n.ChildrenList {
		c.breadthTraverseChildrenList(f)
	}
}

// BreadthFirstTraverseChildrenList
// bread first visit children
func (n *MultiTree) BreadthFirstTraverseChildrenList(f func(*MultiTree) bool) {
	n.breadthTraverseChildrenList(f)
}

// DeepFirstTraverseChildrenList
// deep first traverse children
func (n *MultiTree) DeepFirstTraverseChildrenList(f func(*MultiTree) bool) {
	for _, c := range n.ChildrenList {
		c.deepTraverseChildrenList(f)
	}
}

// TraverseParent
// Traverse a tree node parent
func (n *MultiTree) TraverseParent(f func(*MultiTree) bool) {

	if n.Parent == nil {
		// already is root
		return
	}

	if f(n.Parent) {
		return
	}
	n.Parent.TraverseParent(f)
}

// UnmarshalJSON
// recover from a snapshot,
func (n *MultiTree) UnmarshalJSON(data []byte) error {

	type tmp MultiTree

	t := &tmp{}
	if err := json.Unmarshal(data, t); err != nil {
		return err
	}

	*n = MultiTree(*t)
	n.rebase()
	return nil
}

// rebase
// json.Unmarshal lost origin message, rebase will rich origin info
// must call by root
func (n *MultiTree) rebase() {

	n.Root = n // re pointer to self
	n.rebaseChildrenList()
}

// rebaseChildrenList
// child method for rebase,
// recursive rich tree info
func (n *MultiTree) rebaseChildrenList() {

	for _, c := range n.ChildrenList {

		c.Root = n.Root // pointer to root
		c.Parent = n

		c.rebaseChildrenList()
	}
}
