package queue

import (
	"fmt"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"sync"
	"testing"
	"time"
)

func TestNewQueue(t *testing.T) {

	go func() {
		t.Fatal(http.ListenAndServe(":8080", nil))
	}()

	q := NewQueue(func(i []interface{}) bool {
		if len(i) >= 1000 {
			return true
		}
		return false
	})

	wg := sync.WaitGroup{}

	go func() {

		time.AfterFunc(5*time.Second, func() {

			fmt.Println("i close")
			ol, err := q.Close()
			if err != nil {
				panic(err)
			}
			fmt.Println("i closed", ol)
		})

	}()

	// start the consumer to produce the element
	for i := 0; i < 100000; i++ {

		wg.Add(1)
		go func(i int) {

			defer wg.Done()
			err := q.Append(i)
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Println("sent ", i)
		}(i)
	}

	for i := 0; i < 120000; i++ {

		wg.Add(1)

		go func(i int) {

			defer wg.Done()
			time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
			o, err := q.GET()
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println("got", o)
		}(i)
	}
	wg.Wait()

}
