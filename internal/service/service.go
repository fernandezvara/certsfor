package service

import (
	"context"
	"crypto/x509"
	"crypto/x509/pkix"

	"github.com/fernandezvara/certsfor/internal/structs"
	"github.com/fernandezvara/certsfor/pkg/ca"
	"github.com/fernandezvara/certsfor/pkg/client"
	"github.com/fernandezvara/certsfor/pkg/store"
	"github.com/google/uuid"
)

// Service is the struct used for every request
type Service struct {
	store  store.Store
	client client.Client
	server bool
}

// NewServer creates a Service instance that handles the store directly
func NewServer(store store.Store) *Service {
	return &Service{
		store:  store,
		server: true,
	}
}

// NewClient creates a Server instance that requires a remote server to operate
func NewClient(client client.Client) *Service {
	return &Service{
		client: client,
	}
}

// CACreate is responsible of create a new CA struct with its certificate returning its information
func (s *Service) CACreate(ctx context.Context, subject pkix.Name, years int, months int, days int) (*ca.CA, string, []byte, []byte, error) {

	if s.server {
		return s.caCreateServer(ctx, subject, years, months, days)
	}

	// TODO: call api!
	return nil, "", []byte{}, []byte{}, nil

}

func (s *Service) caCreateServer(ctx context.Context, subject pkix.Name, years int, months int, days int) (*ca.CA, string, []byte, []byte, error) {

	var (
		c           *ca.CA
		cert        []byte
		key         []byte
		id          uuid.UUID
		certificate structs.Certificate
		err         error
	)

	c, cert, key, err = ca.New(subject, years, months, days)
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
func (s *Service) CAGet(caCertificate, caKey []byte) (*ca.CA, error) {

	return ca.FromBytes(caCertificate, caKey)

}

// CertificateGet returns the certificate and its key information
func (s *Service) CertificateGet(ctx context.Context, collection, id string) ([]byte, []byte, error) {

	if s.server {
		return s.certificateGetServer(ctx, collection, id)
	}

	return s.certificateGetClient(ctx, collection, id)

}

func (s *Service) certificateGetServer(ctx context.Context, collection, id string) ([]byte, []byte, error) {

	var (
		value structs.Certificate
		err   error
	)

	err = s.store.Get(ctx, collection, id, &value)

	return value.Certificate, value.Key, err

}

func (s *Service) certificateGetClient(ctx context.Context, collection, id string) ([]byte, []byte, error) {

	// TODO!
	var (
		value structs.Certificate
		err   error
	)

	return value.Certificate, value.Key, err

}

// CertificateCreate creates a new certificate and stores in the store (if server) or POST to the API
func (s *Service) CertificateCreate(ctx context.Context, ca *ca.CA, collection string, info *x509.Certificate) ([]byte, []byte, error) {

	if s.server {
		return s.certificateCreateServer(ctx, ca, collection, info)

	}

	return []byte{}, []byte{}, nil

}

func (s *Service) certificateCreateServer(ctx context.Context, ca *ca.CA, collection string, info *x509.Certificate) ([]byte, []byte, error) {

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
