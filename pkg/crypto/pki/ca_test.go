package pki

import (
	"net"
	"testing"
)

func TestNewCA(t *testing.T) {

	_, err := NewCA("scott-wang.io")
	if err != nil {
		t.Fatal(err)
	}

}

func TestCACreateCSR(t *testing.T) {

	k, _, err := GenerateKeyAndSelfSignedCert("subject")
	if err != nil {
		t.Fatal(err)
	}

	ca, err := NewCA("scott-wang.io")
	if err != nil {
		t.Fatal(err)
	}

	csr, err := ca.CreateCSR(k, "local", []string{"local"}, []net.IP{net.ParseIP("127.0.0.1")})
	if err != nil {
		t.Fatal(err)
	}

	_, err = ca.SignFor(csr, false)
	if err != nil {
		t.Fatal(err)
	}

}
