package main

import (
	"log"
	"runtime"

	"github.com/sxllwx/vulcanus/pkg/scaffold"
	_ "github.com/sxllwx/vulcanus/pkg/scaffold/rest"
	_ "github.com/sxllwx/vulcanus/pkg/scaffold/rest/container"
	_ "github.com/sxllwx/vulcanus/pkg/scaffold/rest/ws"
)

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	if err := scaffold.Cmd().Execute(); err != nil {
		log.Fatal(err)
	}
}
