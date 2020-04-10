package types

import (
	"encoding/json"
	"testing"
)

func TestMultiTree(t *testing.T) {

	root := NewMultiTree("0-0")

	c11 := root.Insert("1-1")
	root.Insert("1-2")

	c21 := c11.Insert("2-1")
	c11.Insert("2-2")

	c21.Insert("3-1")
	c32 := c21.Insert("3-2")

	c32.Insert("4-1")
	last := c32.Insert("4-2")

	breadth, deep := 0, 0

	root.BreadthFirstVisitChildrenList(func(n *MultiTree) {
		t.Logf("breadth-first visit node is (%s)  current-deep (%d) total-deep (%d)", n.Item, n.CurrentDepth, *n.Depth)
		breadth++
	})

	root.DeepFirstVisitChildrenList(func(n *MultiTree) {
		t.Logf("deep-first visit node is (%s)  current-deep (%d) total-deep (%d)", n.Item, n.CurrentDepth, *n.Depth)
		deep++
	})

	if deep != breadth {
		t.Fatal("travese fail")
	}

	last.VisitParent(func(tree *MultiTree) {
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

}
