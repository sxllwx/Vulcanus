package ca

import (
	"github.com/spf13/cobra"
)

// the ca root command
var RootCommand = &cobra.Command{
	Use:  "ca",
	Long: "relate the ca info",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}
