package main

import (
	"runtime"

	"github.com/sxllwx/vulcanus/pkg/scaffold"
	"github.com/sxllwx/vulcanus/pkg/slog"
)

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())
	cmd := scaffold.New()
	if err := cmd.Execute(); err != nil {
		slog.Fatal(err)
	}
}
