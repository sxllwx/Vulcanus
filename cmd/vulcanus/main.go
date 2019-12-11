package main

import (
	"log"
	"runtime"

	"github.com/sxllwx/vulcanus/pkg/scaffold/ca"
	_ "github.com/sxllwx/vulcanus/pkg/scaffold/ca/init"
	"github.com/sxllwx/vulcanus/pkg/scaffold/rest"
	_ "github.com/sxllwx/vulcanus/pkg/scaffold/rest/container"
	_ "github.com/sxllwx/vulcanus/pkg/scaffold/rest/ws"
)

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	rootCommand.AddCommand(rest.RootCommand)
	rootCommand.AddCommand(ca.RootCommand)

	if err := rootCommand.Execute(); err != nil {
		log.Fatal(err)
	}
}
