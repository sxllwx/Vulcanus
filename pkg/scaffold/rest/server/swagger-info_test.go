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
		Package: &Package{
			Name: "server",
		},
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
		Package: &Package{
			Name: "server",
		},
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
			Package: c.Package,
			Model: &Model{
				Name: "Book",
			},
		},
	}

	if err := sg.generate(); err != nil {
		t.Fatal(err)
	}

	if err := wg.generate(); err != nil {
		t.Fatal(err)
	}

	openAPIFD, err := os.OpenFile("open-api.go", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		t.Fatal(err)
	}

	wsFD, err := os.OpenFile("ws.go", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		t.Fatal(err)
	}
	if err := sg.Flush(openAPIFD); err != nil {
		t.Fatal(err)
	}

	if err := wg.Flush(wsFD); err != nil {
		t.Fatal(err)
	}
}
