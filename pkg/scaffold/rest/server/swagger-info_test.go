package server

import (
	"bytes"
	"os"
	"testing"
)

func TestSwagger(t *testing.T) {

	c := swaggerConfig{
		Service: &Service{
			Kind: "books",
		},
		Author: &Author{},
	}

	c.Service.Complete()
	c.Author.Complete()

	sg := swaggerInfoGenerator{
		cache:  &bytes.Buffer{},
		config: &c,
	}

	if err := sg.generate(); err != nil {
		t.Fatal(err)
	}

	t.Log(c.Service.RootURLPrefix)
	t.Logf("%s", sg.cache)
}

func TestE2E(t *testing.T) {

	c := swaggerConfig{
		Service: &Service{
			Kind: "books",
		},
		Author: &Author{},
	}

	c.Service.Complete()
	c.Author.Complete()

	sg := swaggerInfoGenerator{
		cache:  &bytes.Buffer{},
		config: &c,
	}

	wg := webServiceGenerator{
		cache: &bytes.Buffer{},
		config: &webServiceConfig{
			Service: c.Service,
			Model: &Model{
				Name: "Book",
			},
		},
	}

	fh, err := os.OpenFile("demo1.go", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		t.Fatal(err)
	}

	if err := sg.generate(); err != nil {
		t.Fatal(err)
	}

	if err := wg.generate(); err != nil {
		t.Fatal(err)
	}

	if err := sg.Flush(fh); err != nil {
		t.Fatal(err)
	}
	if err := wg.Flush(fh); err != nil {
		t.Fatal(err)
	}
}
