package main

import (
	"runtime"

	"github.com/spf13/cobra"
	"github.com/sxllwx/vulcanus/pkg/scaffold/rest/container"
	"github.com/sxllwx/vulcanus/pkg/scaffold/rest/ws"
)

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	rootCommand := &cobra.Command{
		Use:   "vulcanus",
		Short: "vulcanus is a very awesome golang code generator",
		Run: func(cmd *cobra.Command, args []string) {
			// out put the help
			cmd.Help()
			return
		},
	}

	rootCommand.AddCommand(ws.Command(), container.Command())
	rootCommand.Execute()
}
