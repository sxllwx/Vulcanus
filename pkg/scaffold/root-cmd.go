package scaffold

import (
	"github.com/spf13/cobra"
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

// Cmd
// get the root cmd
func Cmd() *cobra.Command {
	return rootCommand
}
