package net

import (
	"crypto/rand"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func TestReader(t *testing.T) {

	f, err := os.Open("/tmp/rand")
	if err != nil {
		t.Fatal(err)
	}

	rc := DecorateReadCloser(f)
	defer rc.Close()

	ticker := time.NewTicker(100 * time.Millisecond)

	go io.Copy(ioutil.Discard, rc)

	for _ = range ticker.C {
		fmt.Println(rc.BPS())
	}
}

func TestWriter(t *testing.T) {

	test := make([]byte, 1024*1024*1024*4)
	rand.Read(test)

	f, err := os.OpenFile("/tmp/rand", os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		t.Fatal(err)
	}

	wc := DecorateWriteCloser(f)
	defer wc.Close()

	ticker := time.NewTicker(time.Millisecond * 100)
	go wc.Write(test)

	fmt.Println("already start ")
	for _ = range ticker.C {
		fmt.Println(wc.BPS())
	}
}
