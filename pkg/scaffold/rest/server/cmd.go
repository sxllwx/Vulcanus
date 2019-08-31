package server

import (
	"io"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type option struct {

	// src-code package name
	pkg string

	// webservice manage which kind of resource
	kind string
}

func (o *option) run(cmd *cobra.Command, args []string) error {

	s := NewService(o.kind)
	a := NewAuthor("", "", "")
	p := NewPackage(o.pkg)
	m := NewModel(UpperKind(o.kind))

	og := NewOpenAPIGenerator(p, s, a)
	sg := NewWebServiceGenerator(p, s, m)
	return Generate(og, sg)
}

func New() *cobra.Command {

	o := &option{}
	cmd := &cobra.Command{
		Use:   "server",
		Short: "use to generate server code",
		RunE:  o.run,
	}

	cmd.Flags().StringVarP(&o.kind, "kind", "k", "", "your awesome webservice manage the kind of resource")
	cmd.MarkFlagRequired("kind")
	cmd.Flags().StringVarP(&o.pkg, "package", "p", "", "your awesome package")
	cmd.MarkFlagRequired("package")
	return cmd
}

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
			w, err := os.OpenFile(g.SuggestFileName(), os.O_RDWR|os.O_CREATE, 0666)
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
