package scaffold

import "github.com/spf13/cobra"

// scaffold root command
var rootCommand = &cobra.Command{
	Use:   "scaffold",
	Short: "scaffold is a very awesome golang code generator",
	Run: func(cmd *cobra.Command, args []string) {
		// out put the help
		cmd.Help()
		return
	},
}

// New
// get the root cmd
func New() *cobra.Command {
	return rootCommand
}

func AddCommand(cmds ...*cobra.Command) {
	rootCommand.AddCommand(cmds...)
}
