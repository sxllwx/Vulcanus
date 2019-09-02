package ws

import (
	"io/ioutil"
	"testing"
)

func TestWS(t *testing.T) {

	s := NewService("books")
	p := NewPackage("main")
	m := NewModel("Book")
	wsG := NewWebService(p, s, m)

	if err := wsG.Generate(); err != nil {
		t.Fatal(err)
	}

	r, err := ioutil.ReadAll(wsG)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%s", r)
}
