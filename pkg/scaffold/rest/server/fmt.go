package server

import (
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	"golang.org/x/tools/imports"
)

func FormatAndImport(fileName string) error {

	src, err := ioutil.ReadFile(fileName)
	if err != nil {
		return errors.WithMessage(err, "read src")
	}

	result, err := imports.Process("", src, &imports.Options{})
	if err != nil {
		return errors.WithMessage(err, "imports")
	}

	dst, err := os.OpenFile(fileName, os.O_RDWR | os.O_TRUNC, 0666)
	if err != nil{
		return errors.WithMessage(err, "open src")
	}

	defer dst.Close()

	_, err = dst.Write(result)
	if err != nil {
		return errors.WithMessage(err, "write dst")
	}
	return nil
}
