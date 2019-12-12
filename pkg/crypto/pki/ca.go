package pki

import (
	ic "crypto"
	"crypto/x509"
	"fmt"
	"io/ioutil"

	"github.com/juju/errors"
	"github.com/sxllwx/vulcanus/pkg/crypto"
)

type CA struct {
	// private key
	PrivateKey ic.PrivateKey
	// the root cert
	Cert *x509.Certificate
}

func LoadCAFromPEMRSAKeyAndCert(keyPEM []byte, certPEM []byte) (*CA, error) {

	k, err := crypto.LoadRSAPrivateKeyPEM(keyPEM)
	if err != nil {
		return nil, errors.Annotate(err, "load rsa key")
	}

	cert, err := crypto.LoadCertPEM(certPEM)
	if err != nil {
		return nil, errors.Annotate(err, "load cert")
	}

	return &CA{
		PrivateKey: k,
		Cert:       cert,
	}, nil

}

func GenerateCA() (*CA, error) {

	k, err := crypto.GenerateRSAKey()
	if err != nil {
		return nil, errors.Annotate(err, "generate ca key and self cert")
	}

	cert, err := crypto.SelfSignedCert(k)
	if err != nil {
		return nil, errors.Annotate(err, "self signed the cert")
	}

	return &CA{
		PrivateKey: k,
		Cert:       cert,
	}, nil
}

func (ca *CA) FlushALL() {

	var (
		keyFileName  = fmt.Sprintf("%s-key.pem", ca.Cert.Subject.CommonName)
		certFileName = fmt.Sprintf("%s-cert.pem", ca.Cert.Subject.CommonName)
	)

	if err := ioutil.WriteFile(certFileName, crypto.EncodeCertToPEM(ca.Cert), 0644); err != nil {
		panic(err)
	}

	kPEM, err := crypto.EncodePrivateKeyToPEM(ca.PrivateKey)
	if err != nil {
		panic(err)
	}

	if err := ioutil.WriteFile(keyFileName, kPEM, 0644); err != nil {
		panic(err)
	}
}

func (ca *CA) SignFor() {

}
