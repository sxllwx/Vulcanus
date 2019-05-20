package deque

import (
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {

	stop := make(chan struct{})

	q := New(func() {
		fmt.Println("it's persistence op")
	}, stop)

	go func() {

		i := 0
		for {
			if err := q.Insert(i, InsertToHeader); err != nil {
				panic(err)
			}
			i++
		}

	}()

	for {
		o, err := q.Out(OutFromTail)
		if err != nil {
			panic(err)
		}
		fmt.Printf("get %v \n", o)
	}

}
