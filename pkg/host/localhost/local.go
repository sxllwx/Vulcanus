package localhost

import (
	"bytes"
	"os/exec"

	"github.com/juju/errors"
	"github.com/sxllwx/vulcanus/pkg/host"
)

type Client struct {
	cfg *Config
}
type Config struct{}

func NewClient(cfg *Config) host.Interface {

	return &Client{
		cfg: cfg,
	}
}

func (l *Client) Execute(rootCommand string, args ...string) ([]byte, error) {

	buff := &bytes.Buffer{}
	buff.WriteString(rootCommand)
	for _, a := range args {
		buff.WriteString(" ")
		buff.WriteString(a)
	}

	out, err := exec.Command(rootCommand, args...).CombinedOutput()
	if err != nil {
		return nil, errors.Annotate(err, "execute")
	}

	return out, nil
}
