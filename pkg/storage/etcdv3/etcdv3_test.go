package etcdv3

import (
	"bytes"
	"log"
	"os"
	"testing"
	"time"

	"github.com/sxllwx/vulcanus/pkg/storage"
)

var l = log.New(os.Stdout, "test", log.Llongfile)

func newDefaultStorage() (storage.Interface, error) {

	s, err := NewEtcdV3Storage(
		WithHeartbeat(1),
		WithEndpoints("localhost:2379"),
		WithLogger(l),
		WithTimeout(time.Second))

	return s, err
}

func TestNewEtcdV3Storage(t *testing.T) {

	s, err := newDefaultStorage()
	if err != nil {
		t.Fatal(err)
	}
	s.Close()
}

var testCase = []struct {
	k string
	d []byte
}{
	{k: "wx", d: []byte("asd")},
	{k: "wx1", d: []byte("asd1")},
}

func TestClient_PutAndGet(t *testing.T) {

	s, err := newDefaultStorage()
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	defer s.Reset()

	for _, tc := range testCase {
		if err := s.PUT(tc.k, tc.d); err != nil {
			t.Fatal(err)
		}
	}

	for _, tc := range testCase {

		d, err := s.GET(tc.k)
		if err != nil {
			t.Fatal(err)
		}

		if !bytes.Equal(d, tc.d) {
			t.Fatal("not equal")
		}
	}
}

func TestClient_WATCH(t *testing.T) {

	s, err := newDefaultStorage()
	if err != nil {
		t.Fatal(err)
	}

	out, err := s.WATCH("scott")
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		if err := s.PUT("scott", []byte{1, 2, 3, 4}); err != nil {
			t.Fatal(err)
		}
		s.Reset()
		s.Close()
	}()

	for msg := range out {
		t.Logf("get event %#v", msg)
	}
}
