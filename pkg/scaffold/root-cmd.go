package scaffold

import (
	"github.com/spf13/cobra"
	"github.com/sxllwx/vulcanus/pkg/scaffold/rest"
)

// scaffold root command
var rootCommand = &cobra.Command{
	Use:   "vulcanus",
	Short: "vulcanus is a very awesome golang code generator",
	Run: func(cmd *cobra.Command, args []string) {
		// out put the help
		cmd.Help()
		return
	},
}

// New
// get the root cmd
func New() *cobra.Command {
	rootCommand.AddCommand(rest.New())
	return rootCommand
}
