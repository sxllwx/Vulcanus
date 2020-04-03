package container

import (
	"io"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/sxllwx/vulcanus/pkg/scaffold"
	"github.com/sxllwx/vulcanus/pkg/scaffold/rest"
)

type option struct {

	// src-code package name
	pkg string

	// webservice manage which kind of resource
	kind string

	author string
	email  string
	url    string
}

func (o *option) run(cmd *cobra.Command, args []string) error {

	s := rest.NewService(o.kind)
	a := rest.NewAuthor(o.author, o.email, o.url)
	p := rest.NewPackage(o.pkg)
	return scaffold.Generate(NewContainer(p, s, a))
}

func Command() *cobra.Command {

	o := &option{}
	cmd := &cobra.Command{
		Use:   "container",
		Short: "use to generate container code",
		RunE:  o.run,
	}

	cmd.Flags().StringVarP(&o.kind, "kind", "k", "", "your awesome webservice manage the kind of resource")
	cmd.MarkFlagRequired("kind")
	cmd.Flags().StringVarP(&o.pkg, "package", "p", "", "your awesome package")
	cmd.MarkFlagRequired("package")
	cmd.Flags().StringVarP(&o.author, "author", "a", "", "author's name")
	cmd.Flags().StringVarP(&o.email, "email", "e", "", "author's email")
	cmd.Flags().StringVarP(&o.url, "url", "u", "", "author's github url")
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
			if err := scaffold.FormatAndImport(g, g); err != nil {
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
