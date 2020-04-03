package scaffold

import (
	"io"
	"os"

	"github.com/pkg/errors"
)

// Generator
// code generator interface
type Generator interface {
	Generate() error
	io.ReadWriter
	SuggestFileName() string
}

func Generate(gList ...Generator) error {

	for _, g := range gList {

		f := func() error {

			// 1. generate code
			if err := g.Generate(); err != nil {
				return errors.WithMessage(err, "generate src")
			}

			// 2. format and reimport
			if err := FormatAndImport(g, g); err != nil {
				return errors.WithMessage(err, "format and reimport")
			}

			// 3. create file
			w, err := os.OpenFile(g.SuggestFileName(), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
			if err != nil {
				return errors.WithMessagef(err, "open %s", g.SuggestFileName())
			}
			defer w.Close()

			// 4. flush to file
			_, err = io.Copy(w, g)
			if err != nil {

				return errors.WithMessage(err, "copy generated stream to file")
			}
			return nil

		}

		if err := f(); err != nil {
			return err
		}
	}
	return nil
}
