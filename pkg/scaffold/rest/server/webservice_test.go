package server

import (
	"bytes"
	"testing"
)

func TestWS(t *testing.T) {

	a := &Author{}
	s := &Service{
		Kind: "books",
	}

	a.Complete()
	s.Complete()

	wsG := webServiceGenerator{
		cache: &bytes.Buffer{},
		config: &webServiceConfig{
			Service: s,
			Model: &Model{
				Name: "Book",
			},
		},
	}

	if err := wsG.generate(); err != nil {
		t.Fatal(err)
	}

	t.Logf("%s", wsG.cache)
}
