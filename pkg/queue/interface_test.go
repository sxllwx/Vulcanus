package queue

import (
	"context"
	"testing"
	"time"
)

func TestNewQueue(t *testing.T) {

	for i := 0; i < 10000; i++ {
		ctx, _ := context.WithTimeout(context.Background(), time.Second)
		q := New(WithContext(ctx), WithLimiter(notLimit))

		time.Sleep(2 * time.Second)
		if q.Append(1) != ErrAlreadyStopped {
			t.Log(q.Append(1))
			t.Fatal("the ctx no used")
		}
		if _, err := q.Get(); err != ErrAlreadyStopped {
			t.Log(err)
			t.Fatal("the ctx no used")
		}
	}

}
