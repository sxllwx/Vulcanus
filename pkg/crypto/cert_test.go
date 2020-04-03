package crypto

import (
	"io/ioutil"
	"net"
	"testing"
)

func TestSelfSignedCert(t *testing.T) {

	caKey, err := GenerateRSAKey()
	if err != nil {
		t.Fatal(err)
	}

	caCert, err := SelfSignedCert(caKey)
	if err != nil {
		t.Fatal(err)
	}

	caKeyPEM, err := EncodePrivateKeyToPEM(caKey)
	if err != nil {
		t.Fatal(err)
	}

	caCertPEM := EncodeCertToPEM(caCert)

	if err := ioutil.WriteFile("ca-key.pem", caKeyPEM, 0644); err != nil {
		t.Fatal(err)
	}

	if err := ioutil.WriteFile("ca-cert.pem", caCertPEM, 0644); err != nil {
		t.Fatal(err)
	}

	l := net.ParseIP("127.0.0.1")

	subjectKey, err := GenerateRSAKey()
	if err != nil {
		t.Fatal(err)
	}

	subjectKeyPEM, err := EncodePrivateKeyToPEM(subjectKey)
	if err != nil {
		t.Fatal(err)
	}

	if err := ioutil.WriteFile("etcd-key.pem", subjectKeyPEM, 0644); err != nil {
		t.Fatal(err)
	}

	csr, err := CreateCSR("etcd1", []string{"local"}, []net.IP{l}, subjectKey)
	if err != nil {
		t.Fatal(err)
	}

	subjectCert, err := SignForCSR(csr, caCert, caKey, false)
	if err != nil {
		t.Fatal(err)
	}

	subjectCertPEM := EncodeCertToPEM(subjectCert)
	if err := ioutil.WriteFile("etcd-cert.pem", subjectCertPEM, 0644); err != nil {
		t.Fatal(err)
	}
}
