package io

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"sync"
	"testing"
	"time"
)

const (
	testFile = "/tmp/test"
)

func getTestFD(t *testing.T) MeasurableReadWriteCloser {

	f, err := os.OpenFile(testFile, os.O_CREATE|os.O_RDWR|os.O_SYNC, 0666)
	if err != nil {
		t.Fatal(err)
	}
	return DecorateReadWriteCloser(f)
}

func TestReader(t *testing.T) {

	fd := getTestFD(t)

	wwg := sync.WaitGroup{}
	wwg.Add(1)

	// start the write monitor
	go func() {

		ticker := time.NewTicker(time.Second)

		go func() {
			wwg.Wait()
			// write end
			ticker.Stop()
		}()

		// get the real-time metric
		for _ = range ticker.C {
			fmt.Printf("write speed -> %d byte/s\n", fd.WriteMetric().BytesPerSecond())
		}
	}()

	// write loop
	go func() {

		defer wwg.Done()
		test := make([]byte, 1024*8)
		for i := 0; i < 50; i++ {
			rand.Read(test)
			_, err := fd.Write(test)
			if err != nil {
				t.Fatal(err)
			}
		}
	}()
	wwg.Wait()

	// read test start

	rwg := sync.WaitGroup{}
	rwg.Add(1)

	rfd := getTestFD(t)

	// start the write monitor
	go func() {

		ticker := time.NewTicker(time.Second)

		go func() {
			rwg.Wait()
			// write end
			ticker.Stop()
		}()

		// get the real-time metric
		for _ = range ticker.C {
			fmt.Printf("read speed -> %d byte/s\n", rfd.ReadMetric().BytesPerSecond())
		}
	}()

	// read loop
	go func() {

		defer rwg.Done()

		b := make([]byte, 20*3)
		for {
			_, err := rfd.Read(b)

			switch err {
			case nil:
				continue
			case io.EOF:
				return
			default:
				t.Fatal(err)
			}
		}

		fmt.Println("read over")
	}()

	rwg.Wait()

	rfd.Close()

}
