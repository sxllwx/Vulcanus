package scaffold

import (
	"io"
	"io/ioutil"

	"github.com/pkg/errors"
	"golang.org/x/tools/imports"
)

func FormatAndImport(srcR io.Reader, dstW io.Writer) error {

	src, err := ioutil.ReadAll(srcR)
	if err != nil {
		return errors.WithMessage(err, "read src")
	}

	result, err := imports.Process("", src, &imports.Options{
		Comments: true, // keep my comment
	})
	if err != nil {
		return errors.WithMessage(err, "imports")
	}

	_, err = dstW.Write(result)
	if err != nil {
		return errors.WithMessage(err, "write dst")
	}
	return nil
}
