package localhost

import (
	"log"
	"os/exec"

	"github.com/juju/errors"
	"github.com/sxllwx/vulcanus/pkg/host"
)

type Client struct {
	cfg    *Config
	logger *log.Logger
}

type Config struct{}

func NewClient(cfg *Config, l *log.Logger) host.Interface {

	return &Client{
		cfg:    cfg,
		logger: l,
	}
}

func (l *Client) Close() error {
	return nil
}

func (l *Client) Execute(rootCommand string, args ...string) ([]byte, error) {

	out, err := exec.Command(rootCommand, args...).CombinedOutput()
	if err != nil {
		l.logger.Printf("localhost execute (%s, %+v) faild, the err {%s} os output %s", rootCommand, args, err, out)
		return out, errors.Annotatef(err, "run command, os output %s", out)
	}

	l.logger.Printf("localhost execute (%s, %v) success, the os output %s", rootCommand, args, out)
	return out, nil
}
