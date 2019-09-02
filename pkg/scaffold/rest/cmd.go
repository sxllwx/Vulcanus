package rest

import (
	"github.com/spf13/cobra"
	"github.com/sxllwx/vulcanus/pkg/scaffold"
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

// register rest cmd -> scaffold cmd
func init() {
	scaffold.Cmd().AddCommand(restCommand)
}

// Register
// register the rest command to root-cmd
func Register(cmds ...*cobra.Command) {
	restCommand.AddCommand(cmds...)
}
