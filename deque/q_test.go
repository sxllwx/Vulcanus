package deque

import (
	"fmt"
	"k8s.io/apimachinery/pkg/util/rand"
	"testing"
	"time"
)

func TestNew(t *testing.T) {

	stop := make(chan struct{})

	q := New(func() {
		fmt.Println("it's op")
	}, stop)

	go func() {

		for {
			time.Sleep(600 * time.Millisecond)
			if err := q.Insert(rand.String(123), InsertToHeader); err != nil {
				panic(err)
			}
			fmt.Println("pushed a object")
		}

	}()

	for {

		time.Sleep(500 * time.Millisecond)
		o, err := q.Empty()
		if err != nil {
			panic(err)
		}

		fmt.Printf("get %v \n", len(o))
	}

}
