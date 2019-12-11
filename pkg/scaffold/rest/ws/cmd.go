package ws

import (
	"github.com/sxllwx/vulcanus/pkg/scaffold"
	"github.com/sxllwx/vulcanus/pkg/scaffold/rest"

	"github.com/spf13/cobra"
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

// register ws -> rest cmd
func init() {

	o := &option{}
	cmd := &cobra.Command{
		Use:   "ws",
		Short: "use to generate web-service code",
		RunE:  o.run,
	}

	cmd.Flags().StringVarP(&o.kind, "kind", "k", "", "your awesome webservice manage the kind of resource")
	cmd.MarkFlagRequired("kind")
	cmd.Flags().StringVarP(&o.pkg, "package", "p", "", "your awesome package")
	cmd.MarkFlagRequired("package")

	rest.RootCommand.AddCommand(cmd)
	return
}
