package rest

import (
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/juju/errors"
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

func NewClient(endpoint string, versionedPath string) (Interface, error) {

	base, err := url.Parse(endpoint)
	if err != nil {
		return nil, errors.Annotate(err, "parse endpoint")
	}

	tp := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		MaxConnsPerHost:       100,
		IdleConnTimeout:       90 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	c := &http.Client{
		Transport: tp,
	}

	return &RESTClient{
		base:          *base,
		versionedPath: versionedPath,
		c:             c,
	}, nil
}
