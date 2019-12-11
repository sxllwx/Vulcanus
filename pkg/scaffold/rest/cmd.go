package rest

import (
	"github.com/spf13/cobra"
)

//  rest command
var RootCommand = &cobra.Command{
	Use:   "rest",
	Short: "awesome rest golang code generator",
	Run: func(cmd *cobra.Command, args []string) {
		// out put the help
		cmd.Help()
		return
	},
}
