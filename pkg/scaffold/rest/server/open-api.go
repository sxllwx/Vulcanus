package server

import (
	"bytes"
	"text/template"

	"github.com/pkg/errors"
)

type openAPIGenerator struct {
	*bytes.Buffer
	config *openAPIConfig
}

func NewOpenAPIGenerator(p *Package, s *Service, a *Author) Generator {

	return &openAPIGenerator{
		Buffer: &bytes.Buffer{},
		config: &openAPIConfig{
			Package: p,
			Service: s,
			Author:  a,
		},
	}
}

type openAPIConfig struct {
	Package *Package
	Service *Service
	Author  *Author
}

func (g *openAPIGenerator) Generate() error {

	if err := g.generateOpenAPIRegisterFunc(); err != nil {
		return errors.WithMessage(err, "generate open-api register func")
	}

	if err := g.generateRichSwaggerDocFunc(); err != nil {
		return errors.WithMessage(err, "generate rich swagger doc func")
	}

	return nil
}

func (g *openAPIGenerator) generateOpenAPIRegisterFunc() error {

	const tmplt = `package {{.Package.Name}}

import (
	"github.com/emicklei/go-restful-openapi"
	"github.com/go-openapi/spec"
)

func (s *{{.Service.Type}})RegisterOpenAPI(){

	config := restfulspec.Config{
		WebServices: s.container.RegisteredWebServices(), // you control what services are visible
		APIPath:     "/apidocs.json",
		PostBuildSwaggerObjectHandler: s.richSwaggerDoc,
	}
	s.container.Add(restfulspec.NewOpenAPIService(config))
}`

	t, err := template.New("OpenAPITemplate").Parse(tmplt)
	if err != nil {
		return errors.WithMessage(err, "parse template")
	}
	if err := t.Execute(g.Buffer, g.config); err != nil {
		return errors.WithMessage(err, "execute template")
	}
	return nil
}

func (g *openAPIGenerator) generateRichSwaggerDocFunc() error {

	const tmplt = `
func (s *{{.Service.Type}})richSwaggerDoc(swaggerRootDoc *spec.Swagger){

	// TODO: Fix Author Info
	swaggerRootDoc.Info = &spec.Info{
		InfoProps: spec.InfoProps{
			Title:       "{{.Service.Title}}",
			Description: "{{.Service.Description}}",
			Contact: &spec.ContactInfo{
				Name:  "{{.Author.Name}}",
				Email: "{{.Author.Email}}",
				URL:   "{{.Author.URL}}",
			},
			Version: "{{.Service.Version}}",
		},
	}
	swaggerRootDoc.Tags = []spec.Tag{spec.Tag{TagProps: spec.TagProps{
		Name:        "{{.Service.Tag.Name}}",
		Description: "{{.Service.Tag.Description}}",
	}}}
}`

	t, err := template.New("swaggerTemplate").Parse(tmplt)
	if err != nil {
		return errors.WithMessage(err, "parse template")
	}
	if err := t.Execute(g.Buffer, g.config); err != nil {
		return errors.WithMessage(err, "execute template")
	}

	return nil
}

func (g *openAPIGenerator) SuggestFileName() string {
	return openAPISuggestName
}
