package service

import (
	"context"

	"github.com/fernandezvara/certsfor/db/store"
	"github.com/fernandezvara/certsfor/internal/manager"
	"github.com/fernandezvara/certsfor/pkg/client"
	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
)

// Service is the struct used for every request
type Service struct {
	store  store.Store
	client *client.Client
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
func NewAsClient(client *client.Client) *Service {
	return &Service{
		client: client,
	}
}

// Close the store service in a proper way
func (s *Service) Close() error {
	return s.store.Close()
}

// CACreate is responsible of create a new CA struct with its certificate returning its information
func (s *Service) CACreate(ctx context.Context, request client.APICertificateRequest) (*manager.CA, string, []byte, []byte, error) {

	if s.server {
		return s.caCreateServer(ctx, request)
	}

	// TODO: call api!
	return nil, "", []byte{}, []byte{}, nil

}

func (s *Service) caCreateServer(ctx context.Context, request client.APICertificateRequest) (*manager.CA, string, []byte, []byte, error) {

	var (
		c           *manager.CA
		cert        []byte
		key         []byte
		id          uuid.UUID
		certificate client.Certificate
		err         error
	)

	c, cert, key, err = manager.New(request)
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

// CAGet creates a new CA struct from the collection ID
func (s *Service) CAGet(collection string) (*manager.CA, error) {

	var (
		cert client.Certificate
		err  error
	)

	cert, err = s.certificateGetAsServer(context.Background(), collection, "ca")
	if err != nil {
		return nil, err
	}

	return manager.FromBytes(cert.Certificate, cert.Key)

}

// CertificateGet returns the certificate and its key information
func (s *Service) CertificateGet(ctx context.Context, collection, id string) (client.Certificate, error) {

	if s.server {
		return s.certificateGetAsServer(ctx, collection, id)
	}

	return s.certificateGetAsClient(ctx, collection, id)

}

func (s *Service) certificateGetAsServer(ctx context.Context, collection, id string) (certificate client.Certificate, err error) {

	err = s.store.Get(ctx, collection, id, &certificate)
	if err != nil {
		return
	}

	certificate.X509Certificate, err = manager.CertificateFromPEM(certificate.Certificate)

	return

}

func (s *Service) certificateGetAsClient(ctx context.Context, collection, id string) (certificate client.Certificate, err error) {

	// TODO!
	return

}

// CertificateSet creates a new certificate and stores in the store (if server) or POST to the API
func (s *Service) CertificateSet(ctx context.Context, ca *manager.CA, collection string, info client.APICertificateRequest) ([]byte, []byte, error) {

	if s.server {
		return s.certificateSetAsServer(ctx, ca, collection, info)

	}

	return []byte{}, []byte{}, nil

}

// TODO: Make a IsValid for the api request, it must return error if required fields are lost (common name, expirity and key)

func (s *Service) certificateSetAsServer(ctx context.Context, ca *manager.CA, collection string, info client.APICertificateRequest) ([]byte, []byte, error) {

	var (
		certificate client.Certificate
		err         error
	)

	certificate.Certificate, certificate.Key, err = ca.CreateCertificateFromAPI(info)
	if err != nil {
		return []byte{}, []byte{}, err
	}

	err = s.store.Set(ctx, collection, info.DN.CN, certificate)
	if err != nil {
		return []byte{}, []byte{}, err
	}

	return certificate.Certificate, certificate.Key, err
}

// CertificateList returns an array of certificates and its x509 representation
func (s *Service) CertificateList(ctx context.Context, collection string) (certificates []client.Certificate, err error) {

	if s.server {
		return s.certificateListAsServer(ctx, collection)
	}

	return // todo

}

func (s *Service) certificateListAsServer(ctx context.Context, collection string) (certificates []client.Certificate, err error) {

	var (
		mapCertificates []map[string]interface{}
	)

	mapCertificates, err = s.store.GetAll(ctx, collection)
	if err != nil {
		return
	}

	for _, mapCert := range mapCertificates {
		var certificate client.Certificate

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
