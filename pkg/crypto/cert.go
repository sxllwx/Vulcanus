package crypto

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
	"math/big"
	"net"
	"time"

	"github.com/juju/errors"
)

const (
	// the cert TTL
	duration365d = time.Hour * 24 * 365
	// default ca name
	defaultCACommonName = "self-sign-ca"

	PrivateKeyBlockType         = "PRIVATE KEY"
	PublicKeyBlockType          = "PUBLIC KEY"
	CertificateBlockType        = "CERTIFICATE"
	CertificateRequestBlockType = "CERTIFICATE REQUEST"

	ECPrivateKeyBlockType  = "EC PRIVATE KEY"
	RSAPrivateKeyBlockType = "RSA PRIVATE KEY"
)

// SelfSignedCert
// create a self signed cert
func SelfSignedCert(key crypto.PrivateKey) (*x509.Certificate, error) {

	now := time.Now()
	template := x509.Certificate{
		SerialNumber: new(big.Int).SetInt64(0),
		Subject: pkix.Name{
			CommonName: defaultCACommonName,
		},
		NotBefore:             now.UTC(),
		NotAfter:              now.Add(duration365d * 10).UTC(),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		BasicConstraintsValid: true,
		PublicKey:             key.(crypto.Signer).Public(),
		IsCA:                  true,
	}

	ski, err := ComputeSKI(&template)
	if err != nil {
		return nil, errors.Annotate(err, "compute subject key id")
	}

	template.SubjectKeyId = ski

	certDERBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, key.(crypto.Signer).Public(), key)
	if err != nil {
		return nil, errors.Annotate(err, "create cert")
	}

	cert, err := x509.ParseCertificate(certDERBytes)
	if err != nil {
		return nil, errors.Annotatef(err, "parse cert")
	}

	return cert, nil
}

type subjectPublicKeyInfo struct {
	Algorithm        pkix.AlgorithmIdentifier
	SubjectPublicKey asn1.BitString
}

func ComputeSKI(template *x509.Certificate) ([]byte, error) {
	pub := template.PublicKey
	encodedPub, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return nil, err
	}

	var subPKI subjectPublicKeyInfo
	_, err = asn1.Unmarshal(encodedPub, &subPKI)
	if err != nil {
		return nil, err
	}

	pubHash := sha1.Sum(subPKI.SubjectPublicKey.Bytes)
	return pubHash[:], nil
}

// encode the cert to pem
func EncodeCertToPEM(cert *x509.Certificate) []byte {

	return pem.EncodeToMemory(&pem.Block{
		Type:  CertificateBlockType,
		Bytes: cert.Raw,
	})
}

// load the pem bytes to a cert
func LoadCertPEM(raw []byte) (*x509.Certificate, error) {

	block, _ := pem.Decode(raw)
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, errors.Annotate(err, "parse cert")
	}

	return cert, nil
}

func SignForCSR(csr *x509.CertificateRequest, caCert *x509.Certificate, issuerPrivateKey crypto.PrivateKey, isCA bool) (*x509.Certificate, error) {

	var keyUsage x509.KeyUsage

	extKeyUsages := []x509.ExtKeyUsage{}
	if isCA {
		// If the cert is a CA cert, the private key is allowed to sign other certificates.
		keyUsage = x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment | x509.KeyUsageDataEncipherment
		extKeyUsages = append(extKeyUsages, x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth)
	} else {
		// Otherwise the private key is allowed for digital signature and key encipherment.
		keyUsage = x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment
		// For now, we do not differentiate non-CA certs to be used on client auth or server auth.
		extKeyUsages = append(extKeyUsages, x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth)
	}

	now := time.Now()

	serialNumLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNum, err := rand.Int(rand.Reader, serialNumLimit)
	if err != nil {
		return nil, errors.Annotate(err, "rand.Int")
	}

	template := &x509.Certificate{
		SerialNumber: serialNum,

		Subject:            csr.Subject,
		PublicKey:          csr.PublicKey,
		PublicKeyAlgorithm: csr.PublicKeyAlgorithm,
		SignatureAlgorithm: csr.SignatureAlgorithm,

		NotBefore:             now,
		NotAfter:              now.Add(duration365d),
		KeyUsage:              keyUsage,
		ExtKeyUsage:           extKeyUsages,
		IsCA:                  isCA,
		BasicConstraintsValid: true,
		DNSNames:              csr.DNSNames,
		IPAddresses:           csr.IPAddresses,
	}

	certDER, err := x509.CreateCertificate(rand.Reader, template, caCert, csr.PublicKey, issuerPrivateKey)

	cert, err := x509.ParseCertificate(certDER)
	if err != nil {
		return nil, err
	}

	return cert, nil

}

func CreateCSR(subjectCommonName string, dnsSANs []string, ipSANs []net.IP, subjectPrivateKey crypto.PrivateKey) (*x509.CertificateRequest, error) {

	template := x509.CertificateRequest{
		Subject: pkix.Name{
			CommonName: subjectCommonName,
		},
		DNSNames:    dnsSANs,
		IPAddresses: ipSANs,
	}

	template.SignatureAlgorithm = GetSignatureAlgorithm(subjectPrivateKey)

	csrDER, err := x509.CreateCertificateRequest(rand.Reader, &template, subjectPrivateKey)
	if err != nil {
		return nil, errors.Annotate(err, "create certificate request")
	}

	return x509.ParseCertificateRequest(csrDER)
}

func EncodeCSRToPEM(csr *x509.CertificateRequest) []byte {

	return pem.EncodeToMemory(&pem.Block{
		Type:  CertificateRequestBlockType,
		Bytes: csr.Raw,
	})
}

func LoadCSR(raw []byte) (*x509.CertificateRequest, error) {
	block, _ := pem.Decode(raw)
	return x509.ParseCertificateRequest(block.Bytes)
}

func GetSignatureAlgorithm(privateKey crypto.PrivateKey) x509.SignatureAlgorithm {

	switch privateKey.(type) {
	case *rsa.PrivateKey:

		keySize := privateKey.(*rsa.PrivateKey).N.BitLen()
		switch {
		case keySize >= 4096:
			return x509.SHA512WithRSA
		case keySize >= 3072:
			return x509.SHA384WithRSA
		default:
			return x509.SHA256WithRSA
		}
	case *ecdsa.PrivateKey:

		switch privateKey.(*ecdsa.PrivateKey).Curve.Params().BitSize {
		case 256:
			return x509.ECDSAWithSHA256
		case 384:
			return x509.ECDSAWithSHA384
		case 512:
			return x509.ECDSAWithSHA512
		default:
			return x509.UnknownSignatureAlgorithm
		}
	default:
		return x509.UnknownSignatureAlgorithm
	}
}
