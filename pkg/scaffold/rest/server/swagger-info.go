package server

import (
	"bytes"
	"io"
	"text/template"

	"github.com/pkg/errors"
)

type swaggerInfoGenerator struct {
	cache  *bytes.Buffer
	config *swaggerConfig
}

type swaggerConfig struct {
	Package *Package
	Service *Service
	Author  *Author
}

func (g *swaggerInfoGenerator) generate() error {

	if err := g.generateRichSwaggerDocFunc(); err != nil {
		return errors.WithMessage(err, "generate rich swagger doc func")
	}

	if err := g.generateOpenAPIRegisterFunc(); err != nil {
		return errors.WithMessage(err, "generate open-api register func")
	}

	return nil
}

func (g *swaggerInfoGenerator) generateRichSwaggerDocFunc() error {

	const tmplt = `
package {{.Package.Name}}
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
	if err := t.Execute(g.cache, g.config); err != nil {
		return errors.WithMessage(err, "execute template")
	}

	return nil
}

func (g *swaggerInfoGenerator) generateOpenAPIRegisterFunc() error {

	const tmplt = `
func (s *{{.Service.Type}})RegisterOpenAPI(){

	config := restfulspec.Config{
		WebServices: restful.RegisteredWebServices(), // you control what services are visible
		APIPath:     "/apidocs.json",
		PostBuildSwaggerObjectHandler: s.richSwaggerDoc,
	}
	s.container.Add(restfulspec.NewOpenAPIService(config))
}`

	t, err := template.New("OpenAPITemplate").Parse(tmplt)
	if err != nil {
		return errors.WithMessage(err, "parse template")
	}
	if err := t.Execute(g.cache, g.config); err != nil {
		return errors.WithMessage(err, "execute template")
	}
	return nil
}

func (g *swaggerInfoGenerator) Flush(w io.Writer) error {

	_, err := g.cache.WriteTo(w)
	if err != nil {
		return errors.WithMessage(err, "flush to io.writer")
	}
	return nil
}
