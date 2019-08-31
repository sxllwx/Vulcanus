package server

import (
	"bytes"
	"text/template"

	"github.com/pkg/errors"
)

type swaggerInfoGenerator struct {
	cache  *bytes.Buffer
	config *swaggerConfig
}

type swaggerConfig struct {
	Service *Service
	Author  *Author
}

func (g *swaggerInfoGenerator) generate() error {

	const tmplt = `
func (s *{{.Service.Type}})richSwaggerDoc(swaggerRootDoc *spec.Swagger){


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
