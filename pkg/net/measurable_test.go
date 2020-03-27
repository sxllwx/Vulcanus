package net

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestReader(t *testing.T) {

	f, err := os.Open("/tmp/test")
	if err != nil {
		t.Fatal(err)
	}

	rc := DecorateReadCloser(f)
	defer rc.Close()

	ticker := time.NewTicker(time.Second)

	go func() {

		b := make([]byte, 1024*8)
		for i := 0; i < 2000000; i++ {
			rc.Read(b)
		}

	}()

	for _ = range ticker.C {
		fmt.Println(rc.BytesPerSecond())
	}
}

func TestWriter(t *testing.T) {

	f, err := os.OpenFile("/tmp/test_code", os.O_CREATE|os.O_RDWR|os.O_SYNC, 0666)
	if err != nil {
		t.Fatal(err)
	}

	wc := DecorateWriteCloser(f)
	defer wc.Close()

	ticker := time.NewTicker(time.Second)
	go func() {

		test := make([]byte, 1024*1024)
		for i := 0; i < 2000; i++ {
			wc.Write(test)
		}

		fmt.Printf("the write average speed: %v", wc.AverageBytesPerSecond())
		os.Exit(0)
	}()
	for _ = range ticker.C {
		fmt.Println(wc.BytesPerSecond())
	}
}
