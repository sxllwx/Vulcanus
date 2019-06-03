package deque

import (
	"fmt"
	"testing"
	"time"
)

func TestNew(t *testing.T) {

	q := New(func() {
		fmt.Println("it's persistence op")
	})

	time.AfterFunc(2*time.Second, func() {

		q.Stop()
	})

	go func() {

		i := 0
		for {
			if err := q.Insert(i, InsertToHeader); err != nil {
				panic(err)
			}
			i++

			if i == 20 {

				time.Sleep(1 * time.Minute)
			}
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
