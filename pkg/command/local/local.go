package local

import (
	"io"
	"log"
	"os/exec"

	"github.com/sxllwx/vulcanus/pkg/command"
)

type TTY struct {
	cfg    *Config
	logger *log.Logger
}

type Config struct{}

func NewTTY() command.Interface {

	return &TTY{}
}

func (l *TTY) Close() error {
	return nil
}

func (l *TTY) Exec(cmd string, args []string, in io.Reader, out, err io.WriteCloser) error {

	c := exec.Command(cmd, args...)

	c.Stdin = in
	c.Stderr = out
	c.Stdout = err

	return c.Run()
}
