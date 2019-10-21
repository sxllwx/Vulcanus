package localhost

import (
	"testing"
)

var cfg = &Config{}

func TestExecute(t *testing.T) {

	localHost := NewClient(cfg)

	out, err := localHost.Execute("docker")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%s", out)
}
