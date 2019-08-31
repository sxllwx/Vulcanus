package server

import (
	"testing"
)

func TestFormatAndImport(t *testing.T) {

	if err := FormatAndImport("/Users/scott/workspace/go/src/github.com/sxllwx/vulcanus/demo.go"); err != nil{
		t.Fatal(err)
	}
}