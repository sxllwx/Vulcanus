package restclient

import (
	"io"
	"net/http"
	"os"
	"testing"
)

func TestClientGet(t *testing.T) {

	c, err := NewClient("http://localhost:8080", "api/v1.0", http.DefaultTransport)
	if err != nil {
		t.Fatal(err)
	}

	result := c.GET().
		ResourceSet("books").
		Resource("scott").
		Do()

	if err := result.Process(func(response *http.Response) error {
		_, err := io.Copy(os.Stdout, response.Body)
		return err
	}, func(resp *http.Response) error { return nil }); err != nil {
		t.Fatal(err)
	}

}
