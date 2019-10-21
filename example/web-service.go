package main

import (
	"time"

	"github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"
)

// alias the client & server communicate model
// TODO: Fix the struct{} ->  real model
type Books = struct{}

// booksManagerManager
// used to manage resource
type booksManager struct {
	ws *restful.WebService
	// TODO: add other useful field
}

func NewbooksManager() *booksManager {
	s := &booksManager{}
	s.installWebService()
	return s
}

func (s *booksManager) WebService() *restful.WebService {
	return s.ws
}

func (s *booksManager) measureTime(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	now := time.Now()
	chain.ProcessFilter(req, resp)
	time.Now().Sub(now)
}

func (s *booksManager) installWebService() {
	ws := new(restful.WebService)
	ws.
		Path("/apis/v1.0.0/books").
		Consumes(restful.MIME_XML, restful.MIME_JSON).
		Produces(restful.MIME_JSON, restful.MIME_XML)

	tags := []string{"books"}

	ws.Route(ws.POST("").To(s.create).
		// docs
		Doc("create a books").
		Filter(s.measureTime).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Reads(Books{})) // from the request

	ws.Route(ws.PATCH("").To(s.patch).
		// docs
		Doc("patch a books").
		Filter(s.measureTime).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Reads([]byte{})) // from the request

	ws.Route(ws.PUT("/{id}").To(s.update).
		// docs
		Doc("update a books").
		Filter(s.measureTime).
		Param(ws.PathParameter("id", "identifier of the books").DataType("string")).
		// set more rich query condition
		Param(ws.QueryParameter("", "").DataType("")).
		// set more rich header
		Param(ws.HeaderParameter("", "").DataType("")).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Reads(Books{})) // from the request

	ws.Route(ws.GET("/").To(s.list).
		// docs
		Doc("list books").
		// spec a useful filter
		Filter(s.measureTime).
		// spec a spec query condition (the param stay in params)
		Param(ws.QueryParameter("", "").DataType("")).
		// spec a spec query condition (the param stay in header)
		Param(ws.HeaderParameter("", "").DataType("")).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		// the server will provide object-instance for client
		Writes([]Books{}).
		Returns(200, "OK", []Books{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.GET("/{id}").To(s.get).
		// docs
		Doc("get a books").
		// spec a useful filter
		Filter(s.measureTime).
		// spec a spec query condition (the param stay in params)
		Param(ws.PathParameter("id", "identifier of the books").DataType("string")).
		// TODO: QueryParameter
		// TODO: HeaderParameter
		Metadata(restfulspec.KeyOpenAPITags, tags).
		// the server will provide the object-instance
		Writes(Books{}). // on the response
		Returns(200, "OK", Books{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.DELETE("/{id}").To(s.delete).
		// docs
		Doc("delete a books").
		Filter(s.measureTime).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.PathParameter("id", "identifier of the books").DataType("string")))

	s.ws = ws
}

func (s *booksManager) create(request *restful.Request, response *restful.Response) {}
func (s *booksManager) patch(request *restful.Request, response *restful.Response)  {}
func (s *booksManager) list(request *restful.Request, response *restful.Response)   {}
func (s *booksManager) get(request *restful.Request, response *restful.Response)    {}
func (s *booksManager) delete(request *restful.Request, response *restful.Response) {}
func (s *booksManager) update(request *restful.Request, response *restful.Response) {}
