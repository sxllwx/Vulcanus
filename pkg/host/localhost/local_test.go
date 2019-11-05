package localhost

import (
	"log"
	"os"
	"testing"
)

var cfg = &Config{}

func TestExecute(t *testing.T) {

	l := log.New(os.Stdout, "test", log.Lshortfile|log.Ltime)

	localHost := NewClient(cfg, l)

	out, err := localHost.Execute("iptables", "-v")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%s", out)
}
