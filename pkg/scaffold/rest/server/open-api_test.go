package server

import (
	"io/ioutil"
	"testing"
)

func TestOpenAPI(t *testing.T) {

	s := NewService("books")
	p := NewPackage("server")
	a := NewAuthor("", "", "")
	og := NewOpenAPIGenerator(p, s, a)

	if err := og.Generate(); err != nil {
		t.Fatal(err)
	}
	r, err := ioutil.ReadAll(og)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%s", r)
}
