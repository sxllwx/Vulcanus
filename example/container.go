package main

import (
	restful "github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"
	"github.com/go-openapi/spec"
	"net/http"
)

// NewContainer
// create a restful container for hold the web-service
// this container default support for cross origin
func NewContainer() *restful.Container {

	c := restful.NewContainer()
	cors := restful.CrossOriginResourceSharing{
		AllowedHeaders: []string{"Content-Type", "Accept"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		CookiesAllowed: false,
		Container:      c,
	}
	c.Filter(cors.Filter)
	return c
}

// RegisterOpenAPI
// start the open-api docs in container
func RegisterOpenAPI(c *restful.Container) {

	config := restfulspec.Config{
		WebServices:                   c.RegisteredWebServices(), // you control what services are visible
		APIPath:                       "/apidocs.json",
		PostBuildSwaggerObjectHandler: richSwaggerDoc,
	}
	c.Add(restfulspec.NewOpenAPIService(config))
}

// add rich swagger doc, if user need help, he|she can connect with you
func richSwaggerDoc(swaggerRootDoc *spec.Swagger) {

	// TODO: Fix Author Info
	swaggerRootDoc.Info = &spec.Info{
		InfoProps: spec.InfoProps{
			Title:       "BookManagerService",
			Description: "resource for managing book",
			Contact: &spec.ContactInfo{
				Name:  "scott.wang",
				Email: "scottwangsxll@gmail.com",
				URL:   "https://github.com/sxllwx",
			},
			Version: "v1.0.0",
		},
	}
	swaggerRootDoc.Tags = []spec.Tag{spec.Tag{TagProps: spec.TagProps{
		Name:        "book",
		Description: "Managing book",
	}}}
}

// ADD by scott.wang
func main() {

	c := NewContainer()
	NewbookManager(c)
	RegisterOpenAPI(c)

	if err := http.ListenAndServe(":8080", c); err != nil {
		panic(err)
	}

}
