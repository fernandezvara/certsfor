package manager

import (
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
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

func apiTox509Certificate(request structs.APICertificateRequest) (*x509.Certificate, crypto.PrivateKey, error) {

	var (
		cert    x509.Certificate
		key     crypto.PrivateKey
		subject pkix.Name
		err     error
	)

	subject.CommonName = request.DN.CN

	if request.DN.C != "" {
		subject.Country = append(subject.Country, request.DN.C)
	}

	if request.DN.L != "" {
		subject.Locality = append(subject.Locality, request.DN.L)
	}
	if request.DN.O != "" {
		subject.Organization = append(subject.Organization, request.DN.O)
	}
	if request.DN.OU != "" {
		subject.OrganizationalUnit = append(subject.OrganizationalUnit, request.DN.OU)
	}
	if request.DN.P != "" {
		subject.Province = append(subject.Province, request.DN.P)
	}
	if request.DN.PC != "" {
		subject.PostalCode = append(subject.PostalCode, request.DN.PC)
	}
	if request.DN.ST != "" {
		subject.StreetAddress = append(subject.StreetAddress, request.DN.ST)
	}

	cert = x509.Certificate{
		SerialNumber: newSerial(),
		Subject:      subject,
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(time.Duration(request.ExpirationDays*24) * time.Hour),
	}

	// first we need to create the new key that will be used for create any certificates
	switch request.Key {
	case structs.RSA2048:
		key, err = rsa.GenerateKey(rand.Reader, 2048)
	case structs.RSA3072:
		key, err = rsa.GenerateKey(rand.Reader, 3072)
	case structs.RSA4096:
		key, err = rsa.GenerateKey(rand.Reader, 4096)
	case structs.ECDSA224:
		key, err = ecdsa.GenerateKey(elliptic.P224(), rand.Reader)
	case structs.ECDSA256:
		key, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	case structs.ECDSA384:
		key, err = ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	case structs.ECDSA521:
		key, err = ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	default:
		err = ErrKeyInvalid
		return nil, nil, err
	}

	return &cert, key, err

}

// New creates a new CA struct ready to use
func New(request structs.APICertificateRequest) (*CA, []byte, []byte, error) {

	var (
		ca            CA
		caCert, caKey []byte
		err           error
	)

	ca.ca, ca.caKey, err = apiTox509Certificate(request)
	if err != nil {
		return nil, []byte{}, []byte{}, err
	}

	ca.ca.IsCA = true
	ca.ca.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth}
	ca.ca.KeyUsage = x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign
	ca.ca.BasicConstraintsValid = true
	ca.ca.MaxPathLenZero = true

	spkiASN1, err := x509.MarshalPKIXPublicKey(ca.caKey.(crypto.Signer).Public())
	if err != nil {
		return nil, []byte{}, []byte{}, err
	}

	var spki struct {
		Algorithm        pkix.AlgorithmIdentifier
		SubjectPublicKey asn1.BitString
	}
	_, err = asn1.Unmarshal(spkiASN1, &spki)

	skid := sha1.Sum(spki.SubjectPublicKey.Bytes)

	ca.ca.SubjectKeyId = skid[:]

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
		return []byte{}, []byte{}, ErrKeyInvalid
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
	certPrivKeyBytes, err := x509.MarshalPKCS8PrivateKey(c.caKey)
	pem.Encode(certPrivKeyPEM, &pem.Block{
		Type:  FilePrivateKey,
		Bytes: certPrivKeyBytes,
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
	)

	block, _ := pem.Decode(keyPEM)
	if block == nil {
		err = ErrUnparseableFile
		return
	}

	if block.Type != FilePrivateKey {
		err = ErrUnparseableFile
		return
	}

	interfaceKey, err = x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		err = ErrUnparseableFile
		return
	}

	switch interfaceKey.(type) {
	case *rsa.PrivateKey, *ecdsa.PrivateKey:
		// ok
	default:
		err = ErrUnparseableFile
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
