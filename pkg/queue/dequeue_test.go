package queue

import "testing"

func Benchmark_DoubleEndQueue(b *testing.B) {

	for i := 0; i < b.N; i++ {
	}
}

func TestNewDoubleEndQueue(t *testing.T) {

	q := NewDoubleEndQueue()

	c := make(chan struct{})

	go func() {
		o, err := q.PopFromHeader()
		if err != nil {
			t.Fatal(err)
		}
		t.Log(o)
		o, err = q.PopFromHeader()
		if err != nil {
			t.Fatal(err)
		}
		t.Log(o)
		o, err = q.PopFromHeader()
		if err != nil {
			t.Fatal(err)
		}
		t.Log(o)
		o, err = q.PopFromHeader()
		if err != nil {
			t.Fatal(err)
		}
		t.Log(o)

		c <- struct{}{}

	}()
	q.InsertToHeader("1")
	q.InsertToHeader("2")
	q.InsertToTail("3")
	q.InsertToTail("4")

	<-c

}
