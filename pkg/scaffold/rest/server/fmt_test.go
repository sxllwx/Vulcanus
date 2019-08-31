package server

import (
	"io/ioutil"
	"testing"
)

func TestFormatAndImport(t *testing.T) {

	s := NewService("books")
	p := NewPackage("server")
	m := NewModel("Book")
	wsG := NewWebServiceGenerator(p, s, m)

	if err := wsG.Generate(); err != nil {
		t.Fatal(err)
	}

	if err := FormatAndImport(wsG, wsG); err != nil {
		t.Fatal(err)
	}

	r, err := ioutil.ReadAll(wsG)
	if err != nil {
		t.Fatal(err)
	}

	if err := ioutil.WriteFile("demo.go", r, 0666); err != nil {
		t.Fatal(err)
	}

}
