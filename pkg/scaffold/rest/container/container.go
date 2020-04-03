package container

import (
	"bytes"
	"text/template"

	"github.com/pkg/errors"
	"github.com/sxllwx/vulcanus/pkg/scaffold/rest"
)

const (
	containerSuggestName = "container.go"
)

type containerGenerator struct {
	*bytes.Buffer
	config *containerConfig
}

func NewContainer(p rest.Package, s rest.Service, a rest.Author) Generator {

	return &containerGenerator{
		Buffer: &bytes.Buffer{},
		config: &containerConfig{
			Package: p,
			Service: s,
			Author:  a,
		},
	}
}

type containerConfig struct {
	Package rest.Package
	Service rest.Service
	Author  rest.Author
}

func (g *containerGenerator) Generate() error {

	if err := g.generateContainerConstructorFunc(); err != nil {
		return errors.WithMessage(err, "generate container constructor func")
	}
	return nil
}

func (g *containerGenerator) generateContainerConstructorFunc() error {

	const tmplt = `package {{.Package.Name}}

import (
	restful "github.com/emicklei/go-restful"
	"github.com/emicklei/go-restful-openapi"
	"github.com/go-openapi/spec"
)


// NewContainer
// create a restful container for hold the web-service
// this container default support for cross origin
func NewContainer()*restful.Container{

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
func RegisterOpenAPI(c *restful.Container){

	config := restfulspec.Config{
		WebServices: c.RegisteredWebServices(), // you control what services are visible
		APIPath:     "/apidocs.json",
		PostBuildSwaggerObjectHandler: richSwaggerDoc,
	}
	c.Add(restfulspec.NewOpenAPIService(config))
}

// add rich swagger doc, if user need help, he|she can connect with you
func richSwaggerDoc(swaggerRootDoc *spec.Swagger){

	// TODO: Fix Author Info
	swaggerRootDoc.Info = &spec.Info{
		InfoProps: spec.InfoProps{
			Title:       "{{.Service.Title}}",
			Description: "{{.Service.Description}}",
			Contact: &spec.ContactInfo{
				ContactInfoProps: spec.ContactInfoProps{
				    Name:  "{{.Author.Name}}",
				    Email: "{{.Author.Email}}",
				    URL:   "{{.Author.URL}}",
				},
			},
			Version: "{{.Service.Version}}",
		},
	}
	swaggerRootDoc.Tags = []spec.Tag{spec.Tag{TagProps: spec.TagProps{
		Name:        "{{.Service.Tag.Name}}",
		Description: "{{.Service.Tag.Description}}",
	}}}
}
`
	t, err := template.New("NewContainerTemplate").Parse(tmplt)
	if err != nil {
		return errors.WithMessage(err, "parse template")
	}
	if err := t.Execute(g.Buffer, g.config); err != nil {
		return errors.WithMessage(err, "execute template")
	}
	return nil
}

func (g *containerGenerator) SuggestFileName() string {
	return containerSuggestName
}
