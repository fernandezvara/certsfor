package service

import (
	"context"
	"crypto/x509"
	"crypto/x509/pkix"

	"github.com/fernandezvara/certsfor/internal/structs"
	"github.com/fernandezvara/certsfor/pkg/client"
	"github.com/fernandezvara/certsfor/pkg/manager"
	"github.com/fernandezvara/certsfor/pkg/store"
	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
)

// Service is the struct used for every request
type Service struct {
	store  store.Store
	client client.Client
	server bool
}

// NewAsServer creates a Service instance that handles the store directly
func NewAsServer(store store.Store) *Service {
	return &Service{
		store:  store,
		server: true,
	}
}

// NewAsClient creates a Server instance that requires a remote server to operate
func NewAsClient(client client.Client) *Service {
	return &Service{
		client: client,
	}
}

// CACreate is responsible of create a new CA struct with its certificate returning its information
func (s *Service) CACreate(ctx context.Context, subject pkix.Name, years int, months int, days int) (*manager.CA, string, []byte, []byte, error) {

	if s.server {
		return s.caCreateServer(ctx, subject, years, months, days)
	}

	// TODO: call api!
	return nil, "", []byte{}, []byte{}, nil

}

func (s *Service) caCreateServer(ctx context.Context, subject pkix.Name, years int, months int, days int) (*manager.CA, string, []byte, []byte, error) {

	var (
		c           *manager.CA
		cert        []byte
		key         []byte
		id          uuid.UUID
		certificate structs.Certificate
		err         error
	)

	c, cert, key, err = manager.New(subject, years, months, days)
	if err != nil {
		return nil, "", []byte{}, []byte{}, err
	}

	id, err = uuid.NewRandom()
	if err != nil {
		return nil, "", []byte{}, []byte{}, err
	}

	certificate.Certificate = cert
	certificate.Key = key

	err = s.store.Set(ctx, id.String(), "ca", certificate)
	if err != nil {
		return nil, "", []byte{}, []byte{}, err
	}

	return c, id.String(), cert, key, nil

}

// CAGet creates a new CA structs from the cert and key bytes
func (s *Service) CAGet(caCertificate, caKey []byte) (*manager.CA, error) {

	return manager.FromBytes(caCertificate, caKey)

}

// CertificateGet returns the certificate and its key information
func (s *Service) CertificateGet(ctx context.Context, collection, id string) (structs.Certificate, error) {

	if s.server {
		return s.certificateGetAsServer(ctx, collection, id)
	}

	return s.certificateGetAsClient(ctx, collection, id)

}

func (s *Service) certificateGetAsServer(ctx context.Context, collection, id string) (certificate structs.Certificate, err error) {

	err = s.store.Get(ctx, collection, id, &certificate)
	if err != nil {
		return
	}

	certificate.X509Certificate, err = manager.CertificateFromPEM(certificate.Certificate)

	return

}

func (s *Service) certificateGetAsClient(ctx context.Context, collection, id string) (certificate structs.Certificate, err error) {

	// TODO!
	return

}

// CertificateSet creates a new certificate and stores in the store (if server) or POST to the API
func (s *Service) CertificateSet(ctx context.Context, ca *manager.CA, collection string, info *x509.Certificate) ([]byte, []byte, error) {

	if s.server {
		return s.certificateSetAsServer(ctx, ca, collection, info)

	}

	return []byte{}, []byte{}, nil

}

func (s *Service) certificateSetAsServer(ctx context.Context, ca *manager.CA, collection string, info *x509.Certificate) ([]byte, []byte, error) {

	var (
		certificate structs.Certificate
		err         error
	)

	certificate.Certificate, certificate.Key, err = ca.CreateCertificate(info)
	if err != nil {
		return []byte{}, []byte{}, err
	}

	err = s.store.Set(ctx, collection, info.Subject.CommonName, certificate)
	if err != nil {
		return []byte{}, []byte{}, err
	}

	return certificate.Certificate, certificate.Key, err
}

// CertificateList returns an array of certificates and its x509 representation
func (s *Service) CertificateList(ctx context.Context, collection string) (certificates []structs.Certificate, err error) {

	if s.server {
		return s.certificateListAsServer(ctx, collection)
	}

	return // todo

}

func (s *Service) certificateListAsServer(ctx context.Context, collection string) (certificates []structs.Certificate, err error) {

	var (
		mapCertificates []map[string]interface{}
	)

	mapCertificates, err = s.store.GetAll(ctx, collection)
	if err != nil {
		return
	}

	for _, mapCert := range mapCertificates {
		var certificate structs.Certificate

		err = mapstructure.Decode(mapCert, &certificate)
		if err != nil {
			return
		}

		certificate.X509Certificate, err = manager.CertificateFromPEM(certificate.Certificate)
		if err != nil {
			return
		}

		certificates = append(certificates, certificate)
	}

	return
}

// CertificateDelete removes the certificate from the store
func (s *Service) CertificateDelete(ctx context.Context, collection, cn string) (ok bool, err error) {

	if s.server {
		return s.store.Delete(ctx, collection, cn)
	}

	return false, nil // TODO: client

}
