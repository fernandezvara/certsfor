package service_test

import (
	"context"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	_ "github.com/fernandezvara/certsfor/db/badger" // store driver
	"github.com/fernandezvara/certsfor/db/store"
	"github.com/fernandezvara/certsfor/internal/manager"
	"github.com/fernandezvara/certsfor/internal/service"
	"github.com/fernandezvara/certsfor/internal/tests"
	"github.com/fernandezvara/certsfor/pkg/client"
	"github.com/fernandezvara/rest"
	"github.com/stretchr/testify/assert"
)

var (
	caCertificateBytes   []byte
	caKeyBytes           []byte
	certCertificateBytes []byte
	certKeyBytes         []byte
	caID                 string
	caRequest            = client.APICertificateRequest{
		DN: client.APIDN{
			CN: "myca",
			C:  "ES",
			L:  "MyLocality",
			O:  "MyOrganization",
			OU: "MyOU",
			P:  "MyProvince",
			PC: "00000",
			ST: "MyStreet",
		},
		SAN: []string{
			"ca.example.com",
			"192.168.1.1",
		},
		Key:            client.RSA2048,
		ExpirationDays: 90,
		Client:         false,
	}
	certRequest = client.APICertificateRequest{
		DN: client.APIDN{
			CN: "cert",
			C:  "ES",
			L:  "MyLocality",
			O:  "MyOrganization",
			OU: "MyOU",
			P:  "MyProvince",
			PC: "00000",
			ST: "MyStreet",
		},
		SAN: []string{
			"cert.example.com",
			"192.168.1.2",
			"test@example.com",
			"https://an.uri.com/",
		},
		Key:            client.ECDSA521,
		ExpirationDays: 90,
		Client:         false,
	}
)

func TestAPIWithService(t *testing.T) {

	const localIPPort = "127.0.0.1:63998"

	var (
		databaseDir string
		sto         store.Store
		srv         *service.Service
		cli         *client.Client
		srvClient   *service.Service
		status      client.APIStatus
		err         error
	)

	// create temporal directory for the database
	databaseDir, err = ioutil.TempDir("", "cfd")
	assert.Nil(t, err)

	sto, err = store.Open(context.Background(), "badger", databaseDir)
	assert.Nil(t, err)

	srv = service.NewAsServer(sto, tests.Version)
	assert.True(t, srv.Server())

	// must fail, ca not found
	_, err = srv.CAGet("caID-not-found")
	assert.NotNil(t, err)
	assert.Equal(t, rest.ErrNotFound, err)

	// status
	status, err = srv.Status()
	assert.Nil(t, err)
	assert.Equal(t, tests.Version, status.Version)

	testAPI := tests.TestAPI{}

	testAPI.StartAPIWithService(t, localIPPort, "", "", "", 0, false, srv)

	time.Sleep(1 * time.Second)

	cli, err = client.New(localIPPort, "", "", "", false)
	assert.Nil(t, err)

	srvClient = service.NewAsClient(cli, tests.Version)
	assert.False(t, srvClient.Server())
	testStatus(t, srvClient)
	testCreateCA(t, srvClient)
	testCreateCertificate(t, srvClient)
	testGetCertificates(t, srvClient)
	testGetCertificatesParsed(t, srv)
	testListCertificates(t, srvClient)
	testDeleteCertificate(t, srvClient)

	err = testAPI.StopAPI(t)
	assert.Nil(t, err)

	// service without store
	srvUseless := service.NewAsServer(nil, tests.Version)
	assert.Nil(t, srvUseless.Close())

}

func testGetCertificatesParsed(t *testing.T, srv *service.Service) {

	var (
		certificate client.Certificate
		err         error
	)

	certificate, err = srv.CertificateGet(context.Background(), caID, certRequest.DN.CN, 10, true)
	assert.Nil(t, err)
	certificate.X509Certificate, err = manager.CertificateFromPEM(certificate.Certificate)
	assert.Equal(t, certificate.Parsed.Version, certificate.X509Certificate.Version)
	assert.Len(t, certificate.Parsed.IPAddresses, 1)
	assert.Len(t, certificate.Parsed.DNSNames, 1)
	assert.Len(t, certificate.Parsed.EmailAddresses, 1)
	assert.Len(t, certificate.Parsed.URIs, 1)
	assert.Equal(t, certificate.Parsed.NotAfter, certificate.X509Certificate.NotAfter.Unix())
	assert.Equal(t, certificate.Parsed.NotBefore, certificate.X509Certificate.NotBefore.Unix())
	assert.False(t, certificate.Parsed.IsCA)

}

func testStatus(t *testing.T, srv *service.Service) {

	var (
		status client.APIStatus
		err    error
	)

	status, err = srv.Status()
	assert.NotNil(t, status)
	assert.IsType(t, client.APIStatus{}, status)
	assert.Equal(t, tests.Version, status.Version)
	assert.Nil(t, err)

}

func testCreateCA(t *testing.T, srv *service.Service) {

	var (
		ctx context.Context = context.Background()
		err error
	)

	caID, caCertificateBytes, caKeyBytes, err = srv.CACreate(ctx, caRequest)
	assert.Nil(t, err)
	assert.Greater(t, len(caID), 0)
	assert.Greater(t, len(caCertificateBytes), 0)
	assert.Greater(t, len(caKeyBytes), 0)

}

func testCreateCertificate(t *testing.T, srv *service.Service) {

	var (
		ctx           context.Context = context.Background()
		caCertificate []byte
		err           error
	)

	caCertificate, certCertificateBytes, certKeyBytes, err = srv.CertificateSet(ctx, caID, certRequest)
	assert.Nil(t, err)
	assert.Greater(t, len(caID), 0)
	assert.Greater(t, len(caCertificate), 0)
	assert.Greater(t, len(certCertificateBytes), 0)
	assert.Greater(t, len(certKeyBytes), 0)

	assert.Equal(t, caCertificateBytes, caCertificate)

	// must fail, creation of a certificate with the same CN than CA is not allowed
	_, _, _, err = srv.CertificateSet(ctx, caID, caRequest)
	assert.Equal(t, http.StatusText(http.StatusConflict), err.Error()) // client.isError returns http status as text

	// must fail, CA does not exists
	_, _, _, err = srv.CertificateSet(ctx, "ca-not-exists", caRequest)
	assert.Equal(t, http.StatusText(http.StatusNotFound), err.Error())

}

func testGetCertificates(t *testing.T, srv *service.Service) {

	var (
		ctx         context.Context = context.Background()
		certificate client.Certificate
		err         error
	)

	// must fail, ca not found
	_, err = srv.CertificateGet(ctx, "caID not found", "ca", 20, false)
	assert.NotNil(t, err)
	assert.Equal(t, http.StatusText(http.StatusNotFound), err.Error())

	certificate, err = srv.CertificateGet(ctx, caID, "ca", 20, false)
	assert.Nil(t, err)
	assert.Equal(t, caCertificateBytes, certificate.CACertificate)
	assert.Equal(t, caCertificateBytes, certificate.Certificate)
	assert.Equal(t, caKeyBytes, certificate.Key)

	certificate, err = srv.CertificateGet(ctx, caID, certRequest.DN.CN, 20, false)
	assert.Nil(t, err)
	assert.Equal(t, caCertificateBytes, certificate.CACertificate)
	assert.Equal(t, certCertificateBytes, certificate.Certificate)
	assert.Equal(t, certKeyBytes, certificate.Key)

	certificate, err = srv.CertificateGet(ctx, caID, "notfound", 20, false)
	assert.NotNil(t, err)
	assert.Equal(t, http.StatusText(http.StatusNotFound), err.Error())

	certificate, err = srv.CertificateGet(ctx, caID, certRequest.DN.CN, 100, true)
	assert.Nil(t, err)
	assert.Equal(t, caCertificateBytes, certificate.CACertificate)
	assert.NotEqual(t, certCertificateBytes, certificate.Certificate)
	assert.Equal(t, certKeyBytes, certificate.Key)

}

func testListCertificates(t *testing.T, srv *service.Service) {

	var (
		ctx          context.Context = context.Background()
		certificates map[string]client.Certificate
		err          error
	)

	// must fail, ca not found
	_, err = srv.CertificateList(ctx, "caID not found", false)
	assert.NotNil(t, err)
	assert.Equal(t, http.StatusText(http.StatusNotFound), err.Error())

	certificates, err = srv.CertificateList(ctx, caID, false)
	assert.Nil(t, err)
	assert.Len(t, certificates, 2)

}

func testDeleteCertificate(t *testing.T, srv *service.Service) {

	var (
		ctx context.Context = context.Background()
		ok  bool
		err error
	)

	// must fail, ca not found
	ok, err = srv.CertificateDelete(ctx, "caID not found", "id-not-found")
	assert.NotNil(t, err)
	assert.Equal(t, http.StatusText(http.StatusNotFound), err.Error())
	assert.False(t, ok)

	// must fail, common name not found
	ok, err = srv.CertificateDelete(ctx, caID, "id-not-found")
	assert.NotNil(t, err)
	assert.Equal(t, http.StatusText(http.StatusNotFound), err.Error())
	assert.False(t, ok)

	// must fail, ca certificate cannot be deleted
	ok, err = srv.CertificateDelete(ctx, caID, "ca")
	assert.NotNil(t, err)
	assert.Equal(t, http.StatusText(http.StatusConflict), err.Error())
	assert.False(t, ok)

	ok, err = srv.CertificateDelete(ctx, caID, certRequest.DN.CN)
	assert.Nil(t, err)
	assert.True(t, ok)

}
