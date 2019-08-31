package main

import (
	"flag"
	"go/doc"
	"os"

	log "github.com/sxllwx/vulcanus/pkg/slogslog"
)

var (
	goType      = flag.String("go-type", "", "must be set")
	restfulType = flag.String("rest-type", "", "must be set")
)

func main() {

	flag.Parse()

	if len(*goType) == 0 {
		log.Fatal("go-type must be set")
	}
	if len(*restfulType) == 0 {
		log.Fatal("rest-type must be set")
	}

	log.Info(os.Args)

	g := &Generator{}
	err := g.LoadPkg()
	if err != nil {
		log.Fatal(err)
	}

	g.generateTitle()
	g.format()

	for _, e := range g.pkg.TypesInfo.Defs {

		if e == nil {
			continue
		}

		log.Infof("%s", e.Name())
		log.Infof("%#v", e)
	}

	doc.New(g.pkg, ".", 0)
	// log.Infof("%#v", g.pkg.Types.Scope().Lookup("User"))
}
