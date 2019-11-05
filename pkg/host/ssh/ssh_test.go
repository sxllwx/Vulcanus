package ssh

import (
	"log"
	"os"
	"testing"
)

var cfg = &Config{
	Remote:         "192.168.240.101:22",
	User:           "root",
	PrivateKeyFile: "/home/scott/.ssh/id_rsa",
}

var l = log.New(os.Stdout, "test", log.Lshortfile|log.Ltime)

func TestClient_Execute(t *testing.T) {

	c, err := NewClient(cfg, l)
	if err != nil {
		t.Fatal(err)
	}

	out, err := c.Execute("sdasd")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%s", out)
}
