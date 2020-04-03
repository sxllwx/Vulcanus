package ws

import (
	"github.com/spf13/cobra"
	"github.com/sxllwx/vulcanus/pkg/scaffold"
	"github.com/sxllwx/vulcanus/pkg/scaffold/rest"
)

type option struct {

	// src-code package name
	pkg string

	// webservice manage which kind of resource
	kind string
}

func (o *option) run(cmd *cobra.Command, args []string) error {

	s := rest.NewService(o.kind)
	p := rest.NewPackage(o.pkg)
	m := rest.NewModel(rest.UpperKind(o.kind))
	return scaffold.Generate(NewWebService(p, s, m))
}

func Command() *cobra.Command {

	o := &option{}
	cmd := &cobra.Command{
		Use:   "ws",
		Short: "generate web-service code",
		RunE:  o.run,
	}

	cmd.Flags().StringVarP(&o.kind, "kind", "k", "", "resource type")
	cmd.MarkFlagRequired("kind")
	cmd.Flags().StringVarP(&o.pkg, "package", "p", "", "package name")
	cmd.MarkFlagRequired("package")
	return cmd
}
