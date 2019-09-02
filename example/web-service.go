package main

import (
	"time"

	"github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"
)

// alias the client & server communicate model
// TODO: Fix the struct{} ->  real model
type Book = struct{}

// bookManagerManager
// used to manage resource
type bookManager struct {
	ws        *restful.WebService
	container *restful.Container
}

func NewbookManager(c *restful.Container) {

	s := &bookManager{
		container: c,
	}

	s.installWebService()
}

func (s *bookManager) measureTime(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	now := time.Now()
	chain.ProcessFilter(req, resp)
	time.Now().Sub(now)
}

func (s *bookManager) installWebService() {
	ws := new(restful.WebService)
	ws.
		Path("/apis/v1.0.0/book").
		Consumes(restful.MIME_XML, restful.MIME_JSON).
		Produces(restful.MIME_JSON, restful.MIME_XML)

	tags := []string{"book"}

	ws.Route(ws.POST("").To(s.create).
		// docs
		Doc("create a book").
		Filter(s.measureTime).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Reads(Book{})) // from the request

	ws.Route(ws.PATCH("").To(s.patch).
		// docs
		Doc("patch a book").
		Filter(s.measureTime).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Reads([]byte{})) // from the request

	ws.Route(ws.PUT("/{id}").To(s.update).
		// docs
		Doc("update a book").
		Filter(s.measureTime).
		Param(ws.PathParameter("id", "identifier of the book").DataType("string")).
		// set more rich query condition
		Param(ws.QueryParameter("", "").DataType("")).
		// set more rich header
		Param(ws.HeaderParameter("", "").DataType("")).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Reads(Book{})) // from the request

	ws.Route(ws.GET("/").To(s.list).
		// docs
		Doc("list book").
		// spec a useful filter
		Filter(s.measureTime).
		// spec a spec query condition (the param stay in params)
		Param(ws.QueryParameter("", "").DataType("")).
		// spec a spec query condition (the param stay in header)
		Param(ws.HeaderParameter("", "").DataType("")).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		// the server will provide object-instance for client
		Writes([]Book{}).
		Returns(200, "OK", []Book{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.GET("/{id}").To(s.get).
		// docs
		Doc("get a book").
		// spec a useful filter
		Filter(s.measureTime).
		// spec a spec query condition (the param stay in params)
		Param(ws.PathParameter("id", "identifier of the book").DataType("string")).
		// TODO: QueryParameter
		// TODO: HeaderParameter
		Metadata(restfulspec.KeyOpenAPITags, tags).
		// the server will provide the object-instance
		Writes(Book{}). // on the response
		Returns(200, "OK", Book{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.DELETE("/{id}").To(s.delete).
		// docs
		Doc("delete a book").
		Filter(s.measureTime).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.PathParameter("id", "identifier of the book").DataType("string")))

	s.ws = ws
	s.container.Add(ws)
}

func (s *bookManager) create(request *restful.Request, response *restful.Response) {}
func (s *bookManager) patch(request *restful.Request, response *restful.Response)  {}
func (s *bookManager) list(request *restful.Request, response *restful.Response)   {}
func (s *bookManager) get(request *restful.Request, response *restful.Response)    {}
func (s *bookManager) delete(request *restful.Request, response *restful.Response) {}
func (s *bookManager) update(request *restful.Request, response *restful.Response) {}
