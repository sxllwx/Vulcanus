package local

import (
	"os"
	"testing"
)

func TestExecute(t *testing.T) {

	localHost := NewTTY()

	err := localHost.Exec("/bin/ls", []string{}, os.Stdin, os.Stdout, os.Stderr)
	if err != nil {
		t.Fatal(err)
	}
}
