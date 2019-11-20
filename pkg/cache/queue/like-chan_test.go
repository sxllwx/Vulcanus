package queue

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestNewQueue(t *testing.T) {

	wg := sync.WaitGroup{}

	for i := 0; i < 100000; i++ {

		wg.Add(1)

		go func(i int) {

			defer wg.Done()

			ctx, _ := context.WithTimeout(context.Background(), time.Second)
			q := NewLikeChainQueue(WithContext(ctx), WithLimiter(notLimit))
			if err := q.Append("demo"); err != nil {
				t.Fatal(err)
			}
			time.Sleep(2 * time.Second)
			if q.Append(i) != ErrQueueAlreadyStopped {
				t.Log(q.Append(1))
				t.Fatal("the ctx no used")
			}
			if _, err := q.Pop(); err != nil {
				t.Fatal(err)
			}
			if _, err := q.Pop(); err != ErrNoMoreElement {
				t.Fatal(err)
			}

		}(i)
	}

	wg.Wait()
}
