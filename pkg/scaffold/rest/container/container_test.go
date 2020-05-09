package container

import (
	"io/ioutil"
	"testing"

	"github.com/sxllwx/vulcanus/pkg/scaffold/rest"
)

func TestContainer(t *testing.T) {

	s := rest.NewService("books")
	p := rest.NewPackage("main")
	a := rest.NewAuthor("", "", "")
	og := NewContainer(p, s, a)

	if err := og.Generate(); err != nil {
		t.Fatal(err)
	}
	r, err := ioutil.ReadAll(og)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%s", r)
}
