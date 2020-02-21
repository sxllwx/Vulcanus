package tree

import (
	"fmt"
	"testing"
)

func TestTree(t *testing.T) {

	tt := &multiTreeNode{
		obj: 1,
	}

	n2 := tt.Insert(2)
	n3 := tt.Insert(3)
	n3.Insert(6)
	n2.Insert(4)
	n5 := n2.Insert(5)
	n5.Insert(7)
	n5.Insert(8)

	var result []interface{}

	tt.DRL(func(node MultiTreeNode) {
		result = append(result, node.Object())
	})

	fmt.Println(result)

	result = []interface{}{}

	tt.LRD(func(node MultiTreeNode) {
		result = append(result, node.Object())
	})

	fmt.Println(result)
}
