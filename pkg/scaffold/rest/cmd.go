package rest

import (
	"github.com/spf13/cobra"
	"github.com/sxllwx/vulcanus/pkg/scaffold/rest/server"
)

//  rest command
var restCommand = &cobra.Command{
	Use:   "rest",
	Short: "awesome rest golang code generator",
	Run: func(cmd *cobra.Command, args []string) {
		// out put the help
		cmd.Help()
		return
	},
}

// New
// get the rest cmd
func New() *cobra.Command {
	restCommand.AddCommand(server.New())
	return restCommand
}
