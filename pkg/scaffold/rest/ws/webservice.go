package ws

import (
	"bytes"
	"text/template"

	"github.com/sxllwx/vulcanus/pkg/scaffold"

	"github.com/pkg/errors"
	"github.com/sxllwx/vulcanus/pkg/scaffold/rest"
)

const webServiceSuggestName = "web-service.go"

type webServiceGenerator struct {
	*bytes.Buffer
	config *webServiceConfig
}

type webServiceConfig struct {
	Package rest.Package
	Service rest.Service
	Model   rest.Model
}

func NewWebService(p rest.Package, s rest.Service, m rest.Model) scaffold.Generator {

	return &webServiceGenerator{
		Buffer: &bytes.Buffer{},
		config: &webServiceConfig{
			Package: p,
			Service: s,
			Model:   m,
		},
	}
}

func (g *webServiceGenerator) Generate() error {

	if err := g.generateType(); err != nil {
		return errors.WithMessage(err, "generate type")
	}

	return nil

}

func (g *webServiceGenerator) generateType() error {

	const tmplt = `package {{.Package.Name}}

import (
	"time"

	"github.com/emicklei/go-restful"
	"github.com/emicklei/go-restful-openapi"
)

// alias the client & server communicate model
// TODO: Fix the struct{} ->  real model
type {{.Model.Name}} = struct{}


// {{.Service.Type}}Manager
// used to manage resource
type {{.Service.Type}} struct{
   ws *restful.WebService
   // TODO: add other useful field
}

func New{{.Service.Type}}()*{{.Service.Type}}{
   s := &{{.Service.Type}}{}
   s.installWebService()
   return s
}

func (s *{{.Service.Type}}) WebService()*restful.WebService{
	return s.ws
}

func (s *{{.Service.Type}}) installWebService(){
	ws := new(restful.WebService)
	ws.
		Path("{{.Service.RootURLPrefix}}")

	tags := []string{"{{.Service.Tag.Name}}"}

	ws.Route(ws.POST("").To(s.create).
		// docs
		Doc("create a {{.Service.Kind}}").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Reads({{.Model.Name}}{})) // from the request

	ws.Route(ws.PATCH("").To(s.patch).
		// docs
		Doc("patch a {{.Service.Kind}}").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Reads([]byte{})) // from the request

	ws.Route(ws.PUT("/{id}").To(s.update).
		// docs
		Doc("update a {{.Service.Kind}}").
		Param(ws.PathParameter("id", "identifier of the {{.Service.Kind}}").DataType("string")).
		// set more rich query condition
		Param(ws.QueryParameter("", "").DataType("")).
		// set more rich header 
		Param(ws.HeaderParameter("", "").DataType("")).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Reads({{.Model.Name}}{})) // from the request

	ws.Route(ws.GET("/").To(s.list).
		// docs
		Doc("list {{.Service.Kind}}").
		// spec a useful filter 
		// spec a spec query condition (the param stay in params)
		Param(ws.QueryParameter("", "").DataType("")).
		// spec a spec query condition (the param stay in header)
		Param(ws.HeaderParameter("", "").DataType("")).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		// the server will provide object-instance for client
		Writes([]{{.Model.Name}}{}).
		Returns(200, "OK", []{{.Model.Name}}{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.GET("/{id}").To(s.get).
		// docs
		Doc("get a {{.Service.Kind}}").
		// spec a useful filter
		// spec a spec query condition (the param stay in params)
		Param(ws.PathParameter("id", "identifier of the {{.Service.Kind}}").DataType("string")).
		// TODO: QueryParameter 
		// TODO: HeaderParameter 
		Metadata(restfulspec.KeyOpenAPITags, tags).
		// the server will provide the object-instance
		Writes({{.Model.Name}}{}). // on the response
		Returns(200, "OK", {{.Model.Name}}{}).
		Returns(404, "Not Found", nil))



	ws.Route(ws.DELETE("/{id}").To(s.delete).
		// docs
		Doc("delete a {{.Service.Kind}}").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.PathParameter("id", "identifier of the {{.Service.Kind}}").DataType("string")))

	s.ws = ws
}

func (s *{{.Service.Type}})create(request *restful.Request, response *restful.Response){}
func (s *{{.Service.Type}})patch(request *restful.Request, response *restful.Response){}
func (s *{{.Service.Type}})list(request *restful.Request, response *restful.Response){}
func (s *{{.Service.Type}})get(request *restful.Request, response *restful.Response){}
func (s *{{.Service.Type}})delete(request *restful.Request, response *restful.Response){}
func (s *{{.Service.Type}})update(request *restful.Request, response *restful.Response){}

`

	t, err := template.New("types-tplt").Parse(tmplt)
	if err != nil {
		return errors.WithMessage(err, "parse template")
	}

	if err := t.Execute(g.Buffer, g.config); err != nil {
		return errors.WithMessage(err, "execute template")
	}
	return nil
}

func (g *webServiceGenerator) SuggestFileName() string {
	return webServiceSuggestName
}
