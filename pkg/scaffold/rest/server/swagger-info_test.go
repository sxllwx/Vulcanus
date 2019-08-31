package server

import (
	"bytes"
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
