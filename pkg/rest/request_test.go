package rest

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"
)

func TestNewRequest(t *testing.T) {

	baseURL, err := url.Parse("http://localhost:8080")
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.TODO(), 2*time.Second)
	defer cancel()

	r := newRequest(*baseURL, "/api/v1.0", http.MethodGet)
	if err := r.
		Param("name", "scott").
		ResourceSet("books").
		Resource("scott1").
		Context(ctx).
		Do().Process(func(response *http.Response) error {
		io.Copy(os.Stdout, response.Body)
		return nil
	}); err != nil {
		t.Fatal(err)
	}

}

func TestURL(t *testing.T) {

	u, err := url.Parse("tcp://localhost:8080")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%#v", u)
}

func TestURLQuery(t *testing.T) {

	u, err := url.Parse("http://localhost:8080/?hello=world")
	if err != nil {
		t.Fatal(err)
	}

	currentQ := u.Query()
	currentQ.Add("username", "scott")
	u.RawQuery = currentQ.Encode()
	t.Logf("%s", u)
}

func TestURLUserInfo(t *testing.T) {

	u := &url.URL{
		Scheme: "http",
		Host:   "localhost:8080",
	}

	// like git clone http://scott:password@github.com/sxllwx/xxxx
	u.User = url.UserPassword("scott", "password-123")
	t.Logf("%s", u)
}

var (
	urls = []*url.URL{
		{
			Scheme: "http",
			Opaque: "", // encoded opaque data
			//User:       url.UserPassword("scott", "psw"), // username and password information
			User:       nil,
			Host:       "localhost:8888",           // host or host:port
			Path:       "/apis/v1.0.0/books",       // path (relative paths may omit leading slash)
			RawPath:    "/apis/v1.0.0/books/scott", // encoded path hint (see EscapedPath method)
			ForceQuery: false,                      // append a query ('?') even if RawQuery is empty
			RawQuery:   "",                         // encoded query values, without '?'
			Fragment:   "",                         // fragment for references, without '#'
		},
	}
)

func TestURLString(t *testing.T) {

	for _, u := range urls {

		t.Log("request to ", u.String())
	}

}
