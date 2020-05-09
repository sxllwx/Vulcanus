package init

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"io/ioutil"
	"math/big"
	"time"

	"github.com/spf13/cobra"
	"github.com/sxllwx/vulcanus/pkg/scaffold/ca"

	"github.com/juju/errors"
)

const (
	eCPrivateKeyBlockType = "EC PRIVATE KEY"
	certBlockType         = "CERTIFICATE"
)

type option struct {

	// the init private-key and init-cert file
	privateKeyFile string
	certFile       string

	// the init info
	commonName   string
	organization string

	// private-key
	privateKey    *ecdsa.PrivateKey
	privateKeyPEM []byte

	// cert
	cert    *x509.Certificate
	certPEM []byte
}

func init() {

	var o option

	cmd := &cobra.Command{
		Use:  "init",
		Long: "generate the ca root-cert and root-key",

		RunE: func(cmd *cobra.Command, args []string) error {
			return o.run()
		},
	}

	cmd.Flags().StringVarP(&o.privateKeyFile, "private-key-file", "p", "ca-key.pem", "the ca private key file name")
	cmd.Flags().StringVarP(&o.certFile, "cert-file", "c", "ca-cert.pem", "the ca cert file name")
	cmd.Flags().StringVarP(&o.commonName, "common-name", "x", "scott-wang.io", "the ca common name")
	cmd.Flags().StringVarP(&o.organization, "organization", "o", "vulcanus", "the ca organization name")

	ca.RootCommand.AddCommand(cmd)
}

func (o *option) generatePrivateKey() error {

	k, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return errors.Annotate(err, "generate p256 private key")
	}

	derBytes, err := x509.MarshalECPrivateKey(k)
	if err != nil {
		return errors.Annotate(err, "marshal ecdsa private key")
	}

	body := pem.EncodeToMemory(&pem.Block{
		Type:  eCPrivateKeyBlockType,
		Bytes: derBytes,
	})

	o.privateKey = k
	o.privateKeyPEM = body

	return nil
}

func (o *option) selfSignedCert() error {

	now := time.Now()
	tmpl := x509.Certificate{
		SerialNumber: new(big.Int).SetInt64(0),
		Subject: pkix.Name{
			CommonName:   o.commonName,
			Organization: []string{o.organization},
		},
		NotBefore:             now.UTC(),
		NotAfter:              now.Add(duration365d * 10).UTC(),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
		IsCA: true,
	}

	certDERBytes, err := x509.CreateCertificate(rand.Reader, &tmpl, &tmpl, o.privateKey.Public(), o.privateKey)
	if err != nil {
		return errors.Annotate(err, "create cert")
	}

	cert, err := x509.ParseCertificate(certDERBytes)
	if err != nil {
		return errors.Annotatef(err, "parse cert")
	}

	body := pem.EncodeToMemory(&pem.Block{Type: certBlockType, Bytes: certDERBytes})

	o.cert = cert
	o.certPEM = body
	return nil
}

func (o *option) run() error {

	if err := o.generatePrivateKey(); err != nil {
		return errors.Annotate(err, "generate private key")
	}

	if err := o.selfSignedCert(); err != nil {
		return errors.Annotate(err, "self sign cert")
	}

	if err := ioutil.WriteFile(o.privateKeyFile, o.privateKeyPEM, 0644); err != nil {
		return errors.Annotate(err, "flush to private key file")
	}

	if err := ioutil.WriteFile(o.certFile, o.certPEM, 0644); err != nil {
		return errors.Annotate(err, "flush to cert file")
	}

	return nil
}
