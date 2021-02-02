package service_test

import (
	"context"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	_ "github.com/fernandezvara/certsfor/db/badger" // store driver
	"github.com/fernandezvara/certsfor/db/store"
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

	// must fail, ca not found
	_, err = srv.CAGet("caID-not-found")
	assert.NotNil(t, err)
	assert.Equal(t, rest.ErrNotFound, err)

	// status
	status, err = srv.Status()
	assert.Nil(t, err)
	assert.Equal(t, tests.Version, status.Version)

	testAPI := tests.TestAPI{}

	testAPI.StartAPIWithService(t, localIPPort, []byte{}, []byte{}, []byte{}, srv)

	time.Sleep(1 * time.Second)

	cli, err = client.New(localIPPort, "", "", "", false)
	assert.Nil(t, err)

	srvClient = service.NewAsClient(cli, tests.Version)

	testStatus(t, srvClient)
	testCreateCA(t, srvClient)
	testCreateCertificate(t, srvClient)
	testGetCertificates(t, srvClient)
	testListCertificates(t, srvClient)

	err = testAPI.StopAPI(t)
	assert.Nil(t, err)

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
	_, err = srv.CertificateGet(ctx, "caID not found", "ca", 20)
	assert.NotNil(t, err)
	assert.Equal(t, http.StatusText(http.StatusNotFound), err.Error())

	certificate, err = srv.CertificateGet(ctx, caID, "ca", 20)
	assert.Nil(t, err)
	assert.Equal(t, caCertificateBytes, certificate.CACertificate)
	assert.Equal(t, caCertificateBytes, certificate.Certificate)
	assert.Equal(t, caKeyBytes, certificate.Key)

	certificate, err = srv.CertificateGet(ctx, caID, certRequest.DN.CN, 20)
	assert.Nil(t, err)
	assert.Equal(t, caCertificateBytes, certificate.CACertificate)
	assert.Equal(t, certCertificateBytes, certificate.Certificate)
	assert.Equal(t, certKeyBytes, certificate.Key)

	certificate, err = srv.CertificateGet(ctx, caID, certRequest.DN.CN, 100)
	assert.Nil(t, err)
	assert.Equal(t, caCertificateBytes, certificate.CACertificate)
	assert.NotEqual(t, certCertificateBytes, certificate.Certificate)
	assert.Equal(t, certKeyBytes, certificate.Key)

	certificate, err = srv.CertificateGet(ctx, caID, "notfound", 20)
	assert.NotNil(t, err)
	assert.Equal(t, http.StatusText(http.StatusNotFound), err.Error())

}

func testListCertificates(t *testing.T, srv *service.Service) {

	var (
		ctx          context.Context = context.Background()
		certificates map[string]client.Certificate
		err          error
	)

	// must fail, ca not found
	_, err = srv.CertificateList(ctx, "caID not found")
	assert.NotNil(t, err)
	assert.Equal(t, http.StatusText(http.StatusNotFound), err.Error())

	certificates, err = srv.CertificateList(ctx, caID)
	assert.Nil(t, err)
	assert.Len(t, certificates, 2)

}
