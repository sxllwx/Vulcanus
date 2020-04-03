package remote

import (
	"os"
	"testing"
)

var cfg = &Config{
	Remote:         "192.168.240.101:22",
	User:           "root",
	PrivateKeyFile: "/home/scott/.ssh/id_rsa",
}

func TestClient_Execute(t *testing.T) {

	c, err := NewClient(cfg)
	if err != nil {
		t.Fatal(err)
	}

	err = c.Exec("ls", []string{}, os.Stdin, os.Stdout, os.Stderr)
	if err != nil {
		t.Fatal(err)
	}
}
