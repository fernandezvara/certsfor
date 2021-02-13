package service

import (
	"context"
	"crypto"
	"crypto/x509"
	"encoding/base64"
	"time"

	"github.com/fernandezvara/certsfor/db/store"
	"github.com/fernandezvara/certsfor/internal/manager"
	"github.com/fernandezvara/certsfor/pkg/client"
	"github.com/fernandezvara/rest"
	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
)

// Service is the struct used for every request
type Service struct {
	store   store.Store
	client  *client.Client
	server  bool
	version string
}

// NewAsServer creates a Service instance that handles the store directly
func NewAsServer(store store.Store, version string) *Service {
	return &Service{
		version: version,
		store:   store,
		server:  true,
	}
}

// NewAsClient creates a Server instance that requires a remote server to operate
func NewAsClient(client *client.Client, version string) *Service {
	return &Service{
		version: version,
		client:  client,
	}
}

// Server returns true when the library is used as server,
func (s *Service) Server() bool {
	return s.server
}

// Close the store service in a proper way
func (s *Service) Close() error {
	if s.store != nil { // is used?
		return s.store.Close()
	}
	return nil
}

// CACreate is responsible of create a new CA struct with its certificate returning its information
func (s *Service) CACreate(ctx context.Context, request client.APICertificateRequest) (string, []byte, []byte, error) {

	if s.server {
		return s.caCreateServer(ctx, request)
	}

	return s.caCreateClient(ctx, request)

}

func (s *Service) caCreateClient(ctx context.Context, request client.APICertificateRequest) (string, []byte, []byte, error) {

	var (
		certificate client.Certificate
		err         error
	)

	certificate, err = s.client.CACreate(request)
	return certificate.CAID, certificate.Certificate, certificate.Key, err

}

func (s *Service) caCreateServer(ctx context.Context, request client.APICertificateRequest) (string, []byte, []byte, error) {

	var (
		cert        []byte
		key         []byte
		id          uuid.UUID
		certificate client.Certificate
		err         error
	)

	cert, key, err = manager.New(request)
	if err != nil {
		return "", []byte{}, []byte{}, err
	}

	id, err = uuid.NewRandom()
	if err != nil {
		return "", []byte{}, []byte{}, err
	}

	certificate.Certificate = cert
	certificate.Key = key
	certificate.Request = request

	err = s.store.Set(ctx, id.String(), "ca", certificate)
	if err != nil {
		return "", []byte{}, []byte{}, err
	}

	return id.String(), cert, key, nil

}

// CAGet creates a new CA struct from the collection ID
func (s *Service) CAGet(collection string) (*manager.CA, error) {

	var (
		cert client.Certificate
		err  error
	)

	cert, err = s.certificateGetAsServer(context.Background(), collection, "ca", 0, false)
	if err != nil {
		return nil, err
	}

	return manager.FromBytes(cert.Certificate, cert.Key)

}

// CertificateGet returns the certificate and its key information
func (s *Service) CertificateGet(ctx context.Context, collection, id string, remaining int, parse bool) (client.Certificate, error) {

	if s.server {
		return s.certificateGetAsServer(ctx, collection, id, remaining, parse)
	}

	return s.certificateGetAsClient(ctx, collection, id, remaining, parse)

}

func (s *Service) certificateGetAsServer(ctx context.Context, collection, id string, remaining int, parse bool) (certificate client.Certificate, err error) {

	var (
		caCertificate client.Certificate
	)

	// get the CA
	err = s.store.Get(ctx, collection, "ca", &caCertificate)
	if err != nil {
		return
	}

	err = s.store.Get(ctx, collection, id, &certificate)
	if err != nil {
		return
	}

	certificate.X509Certificate, err = manager.CertificateFromPEM(certificate.Certificate)
	certificate.CACertificate = caCertificate.Certificate

	if remaining > 0 {
		if s.IsNearToExpire(certificate, remaining) {

			var (
				ca             *manager.CA
				key            crypto.PrivateKey
				newCertificate *x509.Certificate
			)

			ca, err = manager.FromBytes(caCertificate.Certificate, caCertificate.Key)
			if err != nil {
				return
			}

			// get current Key
			key, err = manager.PrivateKeyFromPEM(certificate.Key)
			if err != nil {
				return
			}

			newCertificate = manager.APITox509Certificate(certificate.Request)

			certificate.Certificate, _, err = ca.CreateCertificate(newCertificate, key)
			if err != nil {
				return
			}

			if parse {
				s.parseCertificate(&certificate)
			}

			err = s.store.Set(ctx, collection, id, certificate)

		}

	}

	return

}

func (s *Service) certificateGetAsClient(ctx context.Context, collection, id string, remaining int, parse bool) (certificate client.Certificate, err error) {

	// api must have a ?renew=20 to return the certificate autorenewed in the API!
	certificate, err = s.client.CertificateGet(collection, id, remaining, parse)
	return

}

// IsNearToExpire returns true if certificate is already expired or remaining days are less than (percent/100)
func (s *Service) IsNearToExpire(certificate client.Certificate, percent int) bool {

	var (
		remainingDays    int64
		maxRemainingDays int64
	)

	maxRemainingDays = certificate.Request.ExpirationDays * (int64(percent) / 100)
	remainingDays = int64(certificate.X509Certificate.NotAfter.Sub(time.Now()).Hours()) / 24

	return remainingDays < maxRemainingDays

}

// CertificateSet creates a new certificate and stores in the store (if server) or POST to the API
func (s *Service) CertificateSet(ctx context.Context, collection string, request client.APICertificateRequest) ([]byte, []byte, []byte, error) {

	if s.server {
		return s.certificateSetAsServer(ctx, collection, request)
	}

	return s.certificateSetAsClient(ctx, collection, request)

}

func (s *Service) certificateSetAsClient(ctx context.Context, collection string, request client.APICertificateRequest) ([]byte, []byte, []byte, error) {

	var (
		response client.Certificate
		err      error
	)

	response, err = s.client.CertificateCreate(collection, request.DN.CN, request)
	return response.CACertificate, response.Certificate, response.Key, err

}

func (s *Service) certificateSetAsServer(ctx context.Context, collection string, request client.APICertificateRequest) ([]byte, []byte, []byte, error) {

	var (
		certificate client.Certificate
		ca          *manager.CA
		err         error
	)

	ca, err = s.CAGet(collection)
	if err != nil {
		return []byte{}, []byte{}, []byte{}, err
	}

	if request.DN.CN == ca.CACertificate().Issuer.CommonName {
		return []byte{}, []byte{}, []byte{}, rest.ErrConflict
	}

	certificate.Certificate, certificate.Key, err = ca.CreateCertificateFromAPI(request)
	if err != nil {
		return []byte{}, []byte{}, []byte{}, err
	}

	certificate.Request = request

	err = s.store.Set(ctx, collection, request.DN.CN, certificate)
	if err != nil {
		return []byte{}, []byte{}, []byte{}, err
	}

	return ca.CACertificateBytes(), certificate.Certificate, certificate.Key, err
}

// CertificateList returns an array of certificates and its x509 representation
func (s *Service) CertificateList(ctx context.Context, collection string, parse bool) (certificates map[string]client.Certificate, err error) {

	if s.server {
		return s.certificateListAsServer(ctx, collection, parse)
	}

	return s.client.CertificateList(collection, parse)

}

func (s *Service) certificateListAsServer(ctx context.Context, collection string, parse bool) (certificates map[string]client.Certificate, err error) {

	var (
		mapCertificates []map[string]interface{}
	)

	certificates = make(map[string]client.Certificate)

	mapCertificates, err = s.store.GetAll(ctx, collection)
	if err != nil {
		return
	}

	for _, mapCert := range mapCertificates {

		var certificate client.Certificate

		// data encode as base64 string, but needs to be decoded to []byte
		if val, ok := mapCert["certificate"].(string); ok {
			mapCert["certificate"], err = base64.StdEncoding.DecodeString(val)
			if err != nil {
				return
			}
		}

		if val, ok := mapCert["key"].(string); ok {
			mapCert["key"], err = base64.StdEncoding.DecodeString(val)
			if err != nil {
				return
			}
		}

		err = mapstructure.Decode(mapCert, &certificate)
		if err != nil {
			return
		}

		certificate.X509Certificate, err = manager.CertificateFromPEM(certificate.Certificate)
		if err != nil {
			return
		}

		if parse {
			s.parseCertificate(&certificate)
		}

		certificates[certificate.X509Certificate.Subject.CommonName] = certificate

	}

	return
}

// type ParsedInfo struct {
// 	Version        int      `json:"version"`
// 	SerialNumber   string   `json:"serial_number"`
// 	NotBefore      int64    `json:"not_before"`
// 	NotAfter       int64    `json:"not_after"`
// 	IsCA           bool     `json:"is_ca"`
// 	DNSNames       []string `json:"dns_names"`
// 	EmailAddresses []string `json:"emails"`
// 	IPAddresses    []string `json:"ips"`
// 	URIs           []string `json:"uris"`
// }

func (s Service) parseCertificate(certificate *client.Certificate) {

	var parsedInfo client.ParsedInfo
	parsedInfo.DNSNames = []string{}
	parsedInfo.EmailAddresses = []string{}
	parsedInfo.IPAddresses = []string{}
	parsedInfo.URIs = []string{}

	parsedInfo.DN = certificate.X509Certificate.Subject.ToRDNSequence().String()
	parsedInfo.Version = certificate.X509Certificate.Version
	parsedInfo.SerialNumber = certificate.X509Certificate.SerialNumber.String()
	parsedInfo.NotBefore = certificate.X509Certificate.NotBefore.Unix()
	parsedInfo.NotAfter = certificate.X509Certificate.NotAfter.Unix()
	parsedInfo.IsCA = certificate.X509Certificate.IsCA
	if certificate.X509Certificate.DNSNames != nil {
		parsedInfo.DNSNames = certificate.X509Certificate.DNSNames
	}
	if certificate.X509Certificate.EmailAddresses != nil {
		parsedInfo.EmailAddresses = certificate.X509Certificate.EmailAddresses
	}
	for _, ip := range certificate.X509Certificate.IPAddresses {
		parsedInfo.IPAddresses = append(parsedInfo.IPAddresses, ip.String())
	}
	for _, uri := range certificate.X509Certificate.URIs {
		parsedInfo.URIs = append(parsedInfo.URIs, uri.String())
	}

	certificate.Parsed = parsedInfo
}

// CertificateDelete removes the certificate from the store
func (s *Service) CertificateDelete(ctx context.Context, collection, cn string) (ok bool, err error) {

	if s.server {
		if cn == "ca" {
			err = rest.ErrConflict
			return
		}
		return s.store.Delete(ctx, collection, cn)
	}

	return s.client.CertificateDelete(collection, cn)

}

// Status returns the status for the service. If used as server it will return ok
// For client request, an API call will be done to ensure availability
func (s *Service) Status() (status client.APIStatus, err error) {

	if s.server {
		status.Version = s.version
		return
	}

	status, err = s.client.Status()
	return

}
