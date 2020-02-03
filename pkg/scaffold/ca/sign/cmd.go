package sign

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"

	"github.com/juju/errors"
	"github.com/spf13/cobra"
	"github.com/sxllwx/vulcanus/pkg/scaffold/ca"
)

type option struct {

	// the ca private-key and ca cert
	caPrivateKeyFile string
	caCertFile       string

	caPrivateKey *ecdsa.PrivateKey
	caCert       *x509.Certificate
}

func (o *option) readCAPrivateKey() error {

	if o.caPrivateKeyFile == "" {
		return errors.New("please spec ca private key file")
	}

	body, err := ioutil.ReadFile(o.caPrivateKeyFile)
	if err != nil {
		return errors.Annotate(err, "read private key file")
	}

	block, _ := pem.Decode(body)
	pk, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return errors.Annotate(err, "parse ecdsa private key")
	}

	o.caPrivateKey = pk
	return nil
}

func (o *option) readCACert() error {

	if o.caCertFile == "" {
		return errors.New("please spec ca cert file")
	}

	body, err := ioutil.ReadFile(o.caCertFile)
	if err != nil {
		return errors.Annotate(err, "read cert file")
	}

	block, _ := pem.Decode(body)
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return errors.Annotate(err, "parse cert")
	}

	o.caCert = cert

	return nil
}

func (o *option) sign() error {
	return nil
}
func (o *option) run() error {

	if err := o.readCAPrivateKey(); err != nil {
		return errors.Annotate(err, "read ca private key")
	}

	if err := o.readCACert(); err != nil {
		return errors.Annotate(err, "read ca cert")
	}

	return nil
}

func init() {

	var o option

	cmd := &cobra.Command{
		Use:  "sign",
		Long: "the ca root-cert sign for a com",

		RunE: func(cmd *cobra.Command, args []string) error {
			return o.run()
		},
	}

	cmd.Flags().StringVarP(&o.caPrivateKeyFile, "ca-private-key-file", "p", "ca-key.pem", "the ca private key file name")
	cmd.Flags().StringVarP(&o.caCertFile, "ca-cert-file", "c", "ca-cert.pem", "the ca cert file name")
	ca.RootCommand.AddCommand(cmd)
}
