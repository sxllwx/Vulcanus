package command

import (
	"io"
)

type Interface interface {
	Exec(cmd string, args []string, in io.Reader, out, err io.WriteCloser) error
	io.Closer
}
