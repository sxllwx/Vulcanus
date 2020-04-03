package ws

import (
	"io/ioutil"
	"testing"

	"github.com/sxllwx/vulcanus/pkg/scaffold/rest"
)

func TestWS(t *testing.T) {

	s := rest.NewService("book")
	p := rest.NewPackage("main")
	m := rest.NewModel("Book")
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
