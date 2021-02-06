package api

import (
	"context"
	"fmt"
	"io/ioutil"
	"testing"
	"time"

	"github.com/fernandezvara/certsfor/db/store"
	"github.com/fernandezvara/certsfor/internal/service"
	"github.com/fernandezvara/certsfor/pkg/client"
	"github.com/stretchr/testify/assert"
)

func TestInternalFileCertificates(t *testing.T) {

	var (
		databaseDir                                 string
		certsDir                                    string
		caCertFile, certCertFile, certKeyFile       string
		caCertBytes, certCertBytes, certKeyBytes    []byte
		caCertBytesC, certCertBytesC, certKeyBytesC []byte
		sto                                         store.Store
		srv                                         *service.Service
		ctx                                         context.Context
		caRequest                                   client.APICertificateRequest
		certRequest                                 client.APICertificateRequest
		caID                                        string
		err                                         error
	)

	// create a service to create certificates for the tests
	//
	// create temporal directory for the database
	databaseDir, err = ioutil.TempDir("", "cfd")
	assert.Nil(t, err)

	// create temporal directory for the database
	certsDir, err = ioutil.TempDir("", "certs")
	assert.Nil(t, err)

	caCertFile = fmt.Sprintf("%s/%s", certsDir, "ca-cert.crt")
	certCertFile = fmt.Sprintf("%s/%s", certsDir, "cert.crt")
	certKeyFile = fmt.Sprintf("%s/%s", certsDir, "key.crt")

	sto, err = store.Open(ctx, "badger", databaseDir)
	assert.Nil(t, err)

	srv = service.NewAsServer(sto, "")

	// create a ca
	caRequest = client.APICertificateRequest{
		DN: client.APIDN{
			CN: "ca",
			O:  "org",
		},
		Key:            client.ECDSA521,
		ExpirationDays: 90,
	}

	caID, caCertBytes, _, err = srv.CACreate(ctx, caRequest)
	assert.Nil(t, err)
	assert.NotEmpty(t, caID)

	err = ioutil.WriteFile(caCertFile, caCertBytes, 0666)
	assert.Nil(t, err)

	// create a certificate
	certRequest = client.APICertificateRequest{
		DN: client.APIDN{
			CN: "cert",
			O:  "org",
		},
		Key:            client.ECDSA521,
		ExpirationDays: 90,
	}

	_, certCertBytes, certKeyBytes, err = srv.CertificateSet(ctx, caID, certRequest)
	assert.Nil(t, err)

	err = ioutil.WriteFile(certCertFile, certCertBytes, 0666)
	assert.Nil(t, err)

	err = ioutil.WriteFile(certKeyFile, certKeyBytes, 0666)
	assert.Nil(t, err)

	var (
		remaining int
		// 	cert, key, cacert                 []byte
		startScheduler bool
	)

	// tests files not found
	_, _, _, _, err = getCertificates("/tmp/non-exist.txt", "", "", remaining, srv)
	assert.Error(t, err)

	_, _, _, _, err = getCertificates(caCertFile, "/tmp/non-exist.txt", "", remaining, srv)
	assert.Error(t, err)

	_, _, _, _, err = getCertificates(caCertFile, certCertFile, "/tmp/non-exist.txt", remaining, srv)
	assert.Error(t, err)

	// test retrival from stores
	//
	// ca does not exits     random uuid for test cfae8b38-57dd-4322-a83f-bc5730689198
	_, _, _, _, err = getCertificates("non-exist", "", "cfae8b38-57dd-4322-a83f-bc5730689198", remaining, srv)
	assert.Error(t, err)

	// certificate does not exitst
	_, _, _, _, err = getCertificates("non-exist", "", caID, remaining, srv)
	assert.Error(t, err)

	// check correctness
	certCertBytesC, certKeyBytesC, caCertBytesC, startScheduler, err = getCertificates("cert", "", caID, 0, srv)
	assert.Nil(t, err)
	assert.Equal(t, certCertBytes, certCertBytesC)
	assert.Equal(t, certKeyBytes, certKeyBytesC)
	assert.Equal(t, caCertBytes, caCertBytesC)
	assert.True(t, startScheduler)

	// func (a *API) Start(apiPort string, tlsCertificate, tlsKey, tlsCaCert string, remaining int,
	//requireClientCertificate bool, outputPaths, errorOutputPaths []string, debug bool) error {

	// start an API with certificates from the database
	var (
		a *API
	)

	a = New(srv, "test")

	go a.Start("127.0.0.1:65000", "cert", "", caID, 100, false, []string{"stdout"}, []string{"stdout"}, true)

	time.Sleep(1 * time.Second)
	err = a.Stop()
	assert.Nil(t, err)

}
