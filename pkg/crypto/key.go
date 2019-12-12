package crypto

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"

	"github.com/juju/errors"
)

// generate a rsa private key
func GenerateRSAKey() (crypto.PrivateKey, error) {

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, errors.Annotate(err, "generate p256 private key")
	}

	return key, nil
}

// generate a ecdsa private key
func GenerateECDSAKey() (crypto.PrivateKey, error) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, errors.Annotate(err, "generate p256 private key")
	}

	return key, nil
}

// encode the private key to PEM
func EncodePrivateKeyToPEM(key crypto.PrivateKey) ([]byte, error) {

	var (
		block = &pem.Block{}
		err   error
	)

	switch key.(type) {

	case *rsa.PrivateKey:

		block.Type = RSAPrivateKeyBlockType
		block.Bytes = x509.MarshalPKCS1PrivateKey(key.(*rsa.PrivateKey))

	case *ecdsa.PrivateKey:

		block.Type = ECPrivateKeyBlockType
		block.Bytes, err = x509.MarshalECPrivateKey(key.(*ecdsa.PrivateKey))
		if err != nil {
			return nil, errors.Annotate(err, "marshal ecdsa private key")
		}

	default:
		return nil, errors.New("unkonow private key typ")
	}

	return pem.EncodeToMemory(block), nil
}

func LoadRSAPrivateKeyPEM(raw []byte) (crypto.PrivateKey, error) {

	block, _ := pem.Decode(raw)
	pk, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, errors.Annotate(err, "parse rsa private key")
	}
	return pk, nil
}

func LoadECDSAPrivateKeyPEM(raw []byte) (crypto.PrivateKey, error) {

	block, _ := pem.Decode(raw)
	pk, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, errors.Annotate(err, "parse ecdsa private key")
	}
	return pk, nil
}
