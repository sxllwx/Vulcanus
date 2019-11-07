package rest

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/juju/errors"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// request
// the model of the http request
type request struct {

	// the underlay http-client
	// the http-client controller the tcp connection between the local && server
	client HTTPClient

	// ---URL---
	u *url.URL
	// the versioned path
	// eg http://localhost:8080/api/v1.0 , the "api/v1.0" is prefix
	versionedPath string
	// the resource set
	// eg http://localhost:8080/api/v1.0/namespaces/scott , the "namespaces" is the resource set
	resourceSet string
	// the resource id
	// eg http://localhost:8080/api/v1.0/namespaces/scott , the "scott" is the resourceID
	resourceID string
	// the request param
	// eg: the book type is cartoon
	param url.Values
	// ---URL---

	// ---HTTP--
	// the method of http
	// POST, GET *
	verb string
	// the http header
	// ContentType or Accept
	header http.Header
	// the request body
	body io.ReadCloser
	// ---HTTP---

	// control the request lifecyle
	ctx context.Context

	// ---MetricHook---
	latencyMetricHook func(verb string, u string, cost time.Duration)
	// ---MetricHook---
}

func (r *request) LatencyHook(hook func(verb string, u string, cost time.Duration)) *request {
	r.latencyMetricHook = hook
	return r
}

func (r *request) ResourceSet(resourceSet string) *request {
	r.resourceSet = resourceSet
	r.u.Path = path.Join(r.u.Path, resourceSet)
	return r
}

func (r *request) Resource(resourceID string) *request {
	r.resourceID = resourceID
	r.u.Path = path.Join(r.u.Path, resourceID)
	return r
}

func (r *request) Header(k string, v ...string) *request {
	r.header[k] = v
	return r
}

func (r *request) Param(k string, v string) *request {
	r.param.Set(k, v)
	return r
}

func (r *request) Body(body io.ReadCloser) *request {
	r.body = body
	return r
}

func (r *request) Context(ctx context.Context) *request {
	r.ctx = ctx
	return r
}

func (r *request) HTTPClient(c HTTPClient) *request {
	r.client = c
	return r
}

// NewRequest
// the base is a bare url, like http://localhost:8080
func newRequest(base url.URL, versionedPath string, verb string) *request {

	base.Path = versionedPath

	return &request{
		versionedPath: versionedPath,
		verb:          verb,
		u:             &base,
		param:         url.Values{},
	}
}

func (r *request) Do() *Result {

	out := &Result{}

	now := time.Now()
	if r.latencyMetricHook != nil {
		defer func() {
			r.latencyMetricHook(r.verb, r.u.String(), time.Since(now))
		}()
	}

	if r.client == nil {
		r.client = http.DefaultClient
	}

	if r.ctx == nil {
		// it's will block forever
		r.ctx = context.Background()
	}

	// add the params
	r.u.RawQuery = r.param.Encode()

	// set header
	request, err := http.NewRequestWithContext(r.ctx, r.verb, r.u.String(), r.body)
	if err != nil {
		out.err = errors.Annotate(err, "new request")
		return out
	}
	request.Header = r.header

	// do
	resp, err := r.client.Do(request)
	if err != nil {
		out.err = errors.Annotate(err, "do request")
		return out
	}

	out.resp = resp
	return out
}

type Result struct {
	// last err, and will cause flow stop
	err  error
	resp *http.Response
}

func (r *Result) Process(handleFunc func(*http.Response) error) error {

	if r.err != nil {
		return r.err
	}

	defer r.resp.Body.Close()

	if err := handleFunc(r.resp); err != nil {
		return errors.Annotate(err, "handle response")
	}
	return nil
}
