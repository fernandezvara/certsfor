package ca

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"time"
)

// CA is the root struct that manages the certificate workflows
type CA struct {
	ca    *x509.Certificate
	caKey *rsa.PrivateKey
}

func (c *CA) newSerial() *big.Int {

	return big.NewInt(time.Now().UnixNano())

}

// New creates a new CA struct ready to use
func New(subject pkix.Name, years, months, days int) (*CA, []byte, []byte, error) {

	var (
		ca            CA
		caCert, caKey []byte
		err           error
	)

	ca.ca = &x509.Certificate{
		SerialNumber:          ca.newSerial(),
		Subject:               subject,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(years, months, days),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	// first we need to create the new key that will be used for create any certificates, also this one
	ca.caKey, err = rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, []byte{}, []byte{}, err
	}

	caCert, caKey, err = ca.CreateCertificate(ca.ca)

	return &ca, caCert, caKey, err

}

// CreateCertificate creates a new certificate from the information passed as request
// returns the PEM files for the certificate and its key
func (c *CA) CreateCertificate(request *x509.Certificate) ([]byte, []byte, error) {

	certBytes, err := x509.CreateCertificate(rand.Reader, request, c.ca, &c.caKey.PublicKey, c.caKey)
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
		Bytes: x509.MarshalPKCS1PrivateKey(c.caKey),
	})

	return certPEM.Bytes(), certPrivKeyPEM.Bytes(), nil
}

// CertificateFromPEM returns a x509.Certificate from the cert PEM and Key bytes
func (c *CA) CertificateFromPEM(certPEM []byte) (cert *x509.Certificate, err error) {

	block, _ := pem.Decode([]byte(certPEM))
	if block == nil {
		return nil, ErrUnparseableFile
	}
	return x509.ParseCertificate(block.Bytes)

}

// PrivateKeyFromPEM returns a rsa.PrivateKey from the PEM bytes
func (c *CA) PrivateKeyFromPEM(keyPEM []byte) (key *rsa.PrivateKey, err error) {

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

	ca.ca, err = ca.CertificateFromPEM(caCertificate)
	if err != nil {
		return nil, err
	}

	ca.caKey, err = ca.PrivateKeyFromPEM(caKey)
	if err != nil {
		return nil, err
	}

	return &ca, err

}
