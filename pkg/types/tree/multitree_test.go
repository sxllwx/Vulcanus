package tree

import (
	"encoding/json"
	"testing"
)

func newNode(id string, name string) interface{} {

	return struct {
		ID   string
		Name string
	}{
		ID:   id,
		Name: name,
	}
}

func TestMultiTree(t *testing.T) {

	root := NewMultiTree(newNode("root", "scott"))

	c11 := root.Insert(newNode("dep-1", "scott-1-1"))
	root.Insert(newNode("dep-1", "scott-1-2"))

	c21 := c11.Insert(newNode("dep-2", "scott-2-1"))
	c11.Insert(newNode("dep-2", "scott-2-2"))

	c21.Insert(newNode("dep-3", "scott-3-1"))
	c32 := c21.Insert(newNode("dep-3", "scott-3-2"))

	c32.Insert(newNode("dep-4", "scott-4-1"))
	last := c32.Insert(newNode("dep-4", "4-2"))

	breadth, deep := 0, 0

	root.BreadthFirstTraverseChildrenList(func(n *MultiTree) {
		t.Logf("breadth-first visit node is (%s)  current-deep (%d) total-deep (%d)", n.Item, n.CurrentDepth, *n.Depth)
		breadth++
	})

	root.DeepFirstTraverseChildrenList(func(n *MultiTree) {
		t.Logf("deep-first visit node is (%s)  current-deep (%d) total-deep (%d)", n.Item, n.CurrentDepth, *n.Depth)
		deep++
	})

	if deep != breadth {
		t.Fatal("traverse fail")
	}

	last.TraverseParent(func(tree *MultiTree) {
		t.Logf("parent (%s) deep (%d) ", tree.Item, tree.CurrentDepth)
	})

	result, err := json.MarshalIndent(root, " ", " ")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%s", result)

	o := NewMultiTree(nil)
	if err := json.Unmarshal(result, o); err != nil {
		t.Fatal(err)
	}

	t.Log(o)

	d := 0
	o.BreadthFirstTraverseChildrenList(func(tree *MultiTree) {
		d++
	})

	if d != deep {
		t.Fatal("unmarshal fail")
	}
}
