package manager

import (
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"time"

	"github.com/fernandezvara/certsfor/internal/structs"
)

// CA is the root struct that manages the certificate workflows
type CA struct {
	ca    *x509.Certificate
	caKey crypto.PrivateKey
}

func newSerial() *big.Int {

	return big.NewInt(time.Now().UnixNano())

}

// // New creates a new CA struct ready to use
// func New(subject pkix.Name, years, months, days int) (*CA, []byte, []byte, error) {

// 	var (
// 		ca            CA
// 		caCert, caKey []byte
// 		err           error
// 	)

// 	ca.ca = &x509.Certificate{
// 		SerialNumber:          newSerial(),
// 		Subject:               subject,
// 		NotBefore:             time.Now(),
// 		NotAfter:              time.Now().AddDate(years, months, days),
// 		IsCA:                  true,
// 		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
// 		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
// 		BasicConstraintsValid: true,
// 	}

// 	// first we need to create the new key that will be used for create any certificates, also this one
// 	ca.caKey, err = rsa.GenerateKey(rand.Reader, 4096)
// 	if err != nil {
// 		return nil, []byte{}, []byte{}, err
// 	}

// 	caCert, caKey, err = ca.CreateCertificate(ca.ca)

// 	return &ca, caCert, caKey, err

// }

func apiTox509Certificate(request structs.APICertificateRequest) *x509.Certificate {

	var (
		cert    x509.Certificate
		subject pkix.Name
	)

	// CN  string   `json:"cn,omitempty"` // common name (required)
	// C   string   `json:"c,omitempty"`  // country
	// L   string   `json:"l,omitempty"`  // locality
	// O   string   `json:"o,omitempty"`  // organization
	// OU  string   `json:"ou,omitempty"` // organization unit
	// P   string   `json:"p,omitempty"`  // province
	// PC  string   `json:"pc,omitempty"` // postal code
	// ST  string   `json:"st,omitempty"` // street

	subject.CommonName = request.CN
	subject.Country = append(subject.Country, request.C)
	subject.Locality = append(subject.Locality, request.L)
	subject.Organization = append(subject.Organization, request.O)
	subject.OrganizationalUnit = append(subject.OrganizationalUnit, request.OU)
	subject.Province = append(subject.Province, request.P)
	subject.PostalCode = append(subject.PostalCode, request.PC)
	subject.StreetAddress = append(subject.StreetAddress, request.ST)

	cert = x509.Certificate{
		SerialNumber: newSerial(),
		Subject:      subject,
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(time.Duration(request.ExpirationDays*24) * time.Hour),
	}

	return &cert

}

// New creates a new CA struct ready to use
func New(request structs.APICertificateRequest) (*CA, []byte, []byte, error) {

	var (
		ca            CA
		caCert, caKey []byte
		err           error
	)

	ca.ca = apiTox509Certificate(request)
	ca.ca.IsCA = true
	ca.ca.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth}
	ca.ca.KeyUsage = x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign
	ca.ca.BasicConstraintsValid = true

	// first we need to create the new key that will be used for create any certificates, also this one

	switch request.Key {
	case structs.RSA2048:
		ca.caKey, err = rsa.GenerateKey(rand.Reader, 2048)
	case structs.RSA4096:
		ca.caKey, err = rsa.GenerateKey(rand.Reader, 4096)
	case structs.ECDSA224:
		ca.caKey, err = ecdsa.GenerateKey(elliptic.P224(), rand.Reader)
	case structs.ECDSA256:
		ca.caKey, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	case structs.ECDSA384:
		ca.caKey, err = ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	case structs.ECDSA521:
		ca.caKey, err = ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	}
	if err != nil {
		return nil, []byte{}, []byte{}, err
	}

	caCert, caKey, err = ca.CreateCertificate(ca.ca)
	if err != nil {
		return nil, []byte{}, []byte{}, err
	}

	return &ca, caCert, caKey, err

}

// CreateCertificate creates a new certificate from the information passed as request
// returns the PEM files for the certificate and its key
func (c *CA) CreateCertificate(request *x509.Certificate) ([]byte, []byte, error) {

	var (
		certBytes []byte
		err       error
	)

	if request.Subject.CommonName == "" {
		return []byte{}, []byte{}, ErrCommonNameBlank
	}

	switch c.caKey.(type) {
	case *rsa.PrivateKey:
		certBytes, err = x509.CreateCertificate(rand.Reader, request, c.ca, &c.caKey.(*rsa.PrivateKey).PublicKey, c.caKey.(*rsa.PrivateKey))
	case *ecdsa.PrivateKey:
		certBytes, err = x509.CreateCertificate(rand.Reader, request, c.ca, &c.caKey.(*ecdsa.PrivateKey).PublicKey, c.caKey.(*ecdsa.PrivateKey))
	default:
		return []byte{}, []byte{}, ErrCAKeyInvalid
	}

	if err != nil {
		return []byte{}, []byte{}, err
	}

	certPEM := new(bytes.Buffer)
	pem.Encode(certPEM, &pem.Block{
		Type:  FileCertificate,
		Bytes: certBytes,
	})

	certPrivKeyPEM := new(bytes.Buffer)
	pem.Encode(certPrivKeyPEM, &pem.Block{
		Type:  FileRSAPrivateKey,
		Bytes: x509.MarshalPKCS1PrivateKey(c.caKey.(*rsa.PrivateKey)),
	})

	return certPEM.Bytes(), certPrivKeyPEM.Bytes(), nil

}

// CertificateFromPEM returns a x509.Certificate from the cert PEM and Key bytes
func CertificateFromPEM(certPEM []byte) (cert *x509.Certificate, err error) {

	block, _ := pem.Decode([]byte(certPEM))
	if block == nil {
		return nil, ErrUnparseableFile
	}
	return x509.ParseCertificate(block.Bytes)

}

// PrivateKeyFromPEM returns a rsa.PrivateKey from the PEM bytes
func PrivateKeyFromPEM(keyPEM []byte) (key *rsa.PrivateKey, err error) {

	var (
		interfaceKey interface{}
		ok           bool
	)

	block, _ := pem.Decode(keyPEM)
	if block == nil {
		err = ErrUnparseableFile
		return
	}

	if block.Type != FileRSAPrivateKey {
		err = ErrUnparseableFile
		return
	}

	interfaceKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		err = ErrUnparseableFile
		return
	}

	key, ok = interfaceKey.(*rsa.PrivateKey)
	if !ok {
		err = ErrUnparseableFile
		return
	}

	return

}

// FromBytes creates a new CA struct from the certificate DER bytes
func FromBytes(caCertificate, caKey []byte) (*CA, error) {

	var (
		ca  CA
		err error
	)

	ca.ca, err = CertificateFromPEM(caCertificate)
	if err != nil {
		return nil, err
	}

	ca.caKey, err = PrivateKeyFromPEM(caKey)
	if err != nil {
		return nil, err
	}

	return &ca, err

}
