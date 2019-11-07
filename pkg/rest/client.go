package rest

import (
	"github.com/juju/errors"
	"net/http"
	"net/url"
)

type RESTClient struct {
	base          url.URL
	versionedPath string
	c             *http.Client
}

func (c *RESTClient) GET() *request {
	return newRequest(c.base, c.versionedPath, http.MethodGet).
		HTTPClient(c.c)
}

func (c *RESTClient) POST() *request {
	return newRequest(c.base, c.versionedPath, http.MethodPost).
		HTTPClient(c.c)
}

func (c *RESTClient) DELETE() *request {
	return newRequest(c.base, c.versionedPath, http.MethodDelete).
		HTTPClient(c.c)
}

func (c *RESTClient) PUT() *request {
	return newRequest(c.base, c.versionedPath, http.MethodPut).
		HTTPClient(c.c)
}

func NewClient(endpoint string, versionedPath string, transport http.RoundTripper) (Interface, error) {

	base, err := url.Parse(endpoint)
	if err != nil {
		return nil, errors.Annotate(err, "parse endpoint")
	}

	if transport == nil {
		transport = http.DefaultTransport
	}

	c := &http.Client{
		Transport: transport,
	}

	return &RESTClient{
		base:          *base,
		versionedPath: versionedPath,
		c:             c,
	}, nil
}
