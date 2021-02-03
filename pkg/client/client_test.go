package client_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"

	_ "github.com/fernandezvara/certsfor/db/badger" // store driver
	"github.com/fernandezvara/certsfor/internal/tests"
	"github.com/fernandezvara/certsfor/pkg/client"
	"github.com/stretchr/testify/assert"
)

const (
	localIPPort    = "127.0.0.1:64002"
	localIPPortSSL = "127.0.0.1:64003"
)

var (
	// all comparations
	caID                   string
	caCertificateBytes     []byte
	certCertificateBytes   []byte
	certKeyBytes           []byte
	clientCertificateBytes []byte
	clientKeyBytes         []byte
	caRequest              client.APICertificateRequest
	certRequest            client.APICertificateRequest
	clientRequest          client.APICertificateRequest
)

func TestClientAPI(t *testing.T) {

	var (
		cli     *client.Client
		testAPI tests.TestAPI
		err     error
	)

	testAPI.StartAPI(t, localIPPort, []byte{}, []byte{}, []byte{}, false) // start HTTP api

	time.Sleep(1 * time.Second)

	cli = createClient(t)
	assert.NotNil(t, cli)

	getStatus(t, cli)
	createCA(t, cli)
	createCertificates(t, cli)
	getCertificate(t, cli)
	listCertificates(t, cli)
	deleteCertificate(t, cli)

	err = testAPI.StopAPI(t)
	assert.Nil(t, err)

	testHTTPSClientAPI(t)

}

func testHTTPSClientAPI(t *testing.T) {

	var (
		certificatesDir string
		certFile        string
		keyFile         string
		caCertFile      string
		cliSSL          *client.Client
		status          client.APIStatus
		testAPI         tests.TestAPI
		err             error
	)

	// create a temporal directory to store certificates for the client
	certificatesDir, err = ioutil.TempDir("", "certs")
	assert.Nil(t, err)

	certFile = fmt.Sprintf("%s/%s", certificatesDir, "cert.crt")
	err = ioutil.WriteFile(certFile, clientCertificateBytes, 0660)
	assert.Nil(t, err)

	keyFile = fmt.Sprintf("%s/%s", certificatesDir, "key.crt")
	err = ioutil.WriteFile(keyFile, clientKeyBytes, 0660)
	assert.Nil(t, err)

	caCertFile = fmt.Sprintf("%s/%s", certificatesDir, "cacert.crt")
	err = ioutil.WriteFile(caCertFile, caCertificateBytes, 0660)
	assert.Nil(t, err)

	// https api
	testAPI.StartAPI(t, localIPPortSSL, certCertificateBytes, certKeyBytes, caCertificateBytes, true)

	time.Sleep(1 * time.Second)

	cliSSL, err = client.NewWithConnectionTimeouts(localIPPortSSL, caCertFile, certFile, keyFile, false, 100*time.Millisecond, 100*time.Millisecond, 500*time.Millisecond)
	assert.Nil(t, err)

	status, err = cliSSL.Status()
	assert.Nil(t, err)
	assert.Equal(t, tests.Version, status.Version)

	// remove certificates directory
	err = os.RemoveAll(certificatesDir)
	assert.Nil(t, err)

	err = testAPI.StopAPI(t)
	assert.Nil(t, err)

}

func TestClientWithErrors(t *testing.T) {

	var (
		cliWithErrors *client.Client
		err           error
	)

	// cannot connect to the remote API
	cliWithErrors, err = client.NewWithConnectionTimeouts("127.0.0.1:64100", "", "", "", false, 100*time.Millisecond, 100*time.Millisecond, 500*time.Millisecond)
	assert.Nil(t, err)

	_, err = cliWithErrors.Status()
	assert.Error(t, err)

	_, err = cliWithErrors.CACreate(client.APICertificateRequest{})
	assert.Error(t, err)

	_, err = cliWithErrors.CertificateCreate("ca-uuid", "common-name", client.APICertificateRequest{})
	assert.Error(t, err)

	_, err = cliWithErrors.CertificateGet("ca-uuid", "common-name", 20)
	assert.Error(t, err)

	_, err = cliWithErrors.CertificateDelete("ca-uuid", "common-name")
	assert.Error(t, err)

	// client with non existent certificates or keys
	// error ca cert
	cliWithErrors, err = client.NewWithConnectionTimeouts("127.0.0.1:64100", "/non-existent/cacert.txt", "", "", false, 100*time.Millisecond, 100*time.Millisecond, 500*time.Millisecond)
	assert.Error(t, err)

	// error cert or key
	cliWithErrors, err = client.NewWithConnectionTimeouts("127.0.0.1:64100", "", "/non-existent/cert.txt", "/non-existent/key.txt", false, 100*time.Millisecond, 100*time.Millisecond, 500*time.Millisecond)
	assert.Error(t, err)

	// use system certificates option
	cliWithErrors, err = client.NewWithConnectionTimeouts("127.0.0.1:64100", "", "", "", true, 100*time.Millisecond, 100*time.Millisecond, 500*time.Millisecond)
	assert.Nil(t, err)

	_, err = cliWithErrors.Status()
	assert.Error(t, err)

}

func createClient(t *testing.T) (cli *client.Client) {

	var err error

	cli, err = client.New(localIPPort, "", "", "", false)
	assert.Nil(t, err)

	return cli

}

func getStatus(t *testing.T, cli *client.Client) {

	var (
		status client.APIStatus
		err    error
	)

	status, err = cli.Status()
	assert.Nil(t, err)
	assert.Equal(t, tests.Version, status.Version)

}

func createCA(t *testing.T, cli *client.Client) {

	var (
		caCertificate client.Certificate
		err           error
	)

	// must fail, request is incomplete
	_, err = cli.CACreate(caRequest)
	assert.Error(t, err)
	assert.Equal(t, http.StatusText(http.StatusBadRequest), err.Error())

	caRequest = client.APICertificateRequest{
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
		Key:            "rsa:4096",
		ExpirationDays: 90,
		Client:         false,
	}

	caCertificate, err = cli.CACreate(caRequest)
	assert.Nil(t, err)
	assert.Equal(t, caRequest.DN.CN, caCertificate.Request.DN.CN)
	assert.Len(t, caCertificate.Request.SAN, 2)
	assert.Equal(t, caCertificate.Certificate, caCertificate.CACertificate)

	caID = caCertificate.CAID
	caCertificateBytes = caCertificate.CACertificate

}

func createCertificates(t *testing.T, cli *client.Client) {

	var (
		certCertificate   client.Certificate
		clientCertificate client.Certificate
		err               error
	)

	certRequest = client.APICertificateRequest{
		DN: client.APIDN{
			CN: "willfail",
		},
	}

	// must fail, request is incomplete
	_, err = cli.CertificateCreate(caID, "willfail", certRequest)
	assert.Error(t, err)
	assert.Equal(t, http.StatusText(http.StatusBadRequest), err.Error())

	certRequest = client.APICertificateRequest{
		DN: client.APIDN{
			CN: "does-not-match",
		},
		Key:            "ecdsa:521",
		ExpirationDays: 90,
	}

	// must fail, common name does not match
	_, err = cli.CertificateCreate(caID, "willfail", certRequest)
	assert.Error(t, err)
	assert.Equal(t, http.StatusText(http.StatusConflict), err.Error())

	certRequest = client.APICertificateRequest{
		DN: client.APIDN{
			CN: "service1",
			C:  "ES",
			L:  "MyLocality",
			O:  "MyOrganization",
			OU: "MyOU",
			P:  "MyProvince",
			PC: "00000",
			ST: "MyStreet",
		},
		SAN: []string{
			"service1.example.com",
			"127.0.0.1",
		},
		Key:            "rsa:4096",
		ExpirationDays: 90,
		Client:         false,
	}

	certCertificate, err = cli.CertificateCreate(caID, certRequest.DN.CN, certRequest)
	assert.Nil(t, err)
	assert.Equal(t, certRequest.DN.CN, certCertificate.Request.DN.CN)
	assert.Len(t, certCertificate.Request.SAN, 2)
	assert.Equal(t, certCertificate.CACertificate, caCertificateBytes)
	certCertificateBytes = certCertificate.Certificate
	certKeyBytes = certCertificate.Key

	clientRequest = client.APICertificateRequest{
		DN: client.APIDN{
			CN: "client",
			C:  "ES",
			L:  "MyLocality",
			O:  "MyOrganization",
			OU: "MyOU",
			P:  "MyProvince",
			PC: "00000",
			ST: "MyStreet",
		},
		SAN: []string{
			"client@example.com",
		},
		Key:            "ecdsa:521",
		ExpirationDays: 90,
		Client:         true,
	}

	clientCertificate, err = cli.CertificateCreate(caID, clientRequest.DN.CN, clientRequest)
	assert.Nil(t, err)
	assert.Equal(t, clientRequest.DN.CN, clientCertificate.Request.DN.CN)
	assert.Len(t, clientCertificate.Request.SAN, 1)
	assert.Equal(t, clientCertificate.CACertificate, caCertificateBytes)
	clientCertificateBytes = clientCertificate.Certificate
	clientKeyBytes = clientCertificate.Key

}

func getCertificate(t *testing.T, cli *client.Client) {

	var (
		certCertificate client.Certificate
		err             error
	)

	// 404 - Not found
	_, err = cli.CertificateGet(caID, "404", 20)
	assert.Error(t, err)
	assert.Equal(t, http.StatusText(http.StatusNotFound), err.Error())

	// 404 - Not found
	_, err = cli.CertificateGet("1234", "404", 20)
	assert.Error(t, err)
	assert.Equal(t, http.StatusText(http.StatusNotFound), err.Error())

	// 200 - OK, must match with the certificated created
	certCertificate, err = cli.CertificateGet(caID, certRequest.DN.CN, 20)
	assert.Nil(t, err)

	assert.Equal(t, certRequest.DN.CN, certCertificate.Request.DN.CN)
	assert.Len(t, certCertificate.Request.SAN, 2)
	assert.Equal(t, certCertificate.CACertificate, caCertificateBytes)
	assert.Equal(t, certCertificate.Certificate, certCertificateBytes)
	assert.Equal(t, certCertificate.Key, certKeyBytes)

	// 200 - OK, must match with the certificated created, but renewed
	certCertificate, err = cli.CertificateGet(caID, certRequest.DN.CN, 100)
	assert.Nil(t, err)

	assert.Equal(t, certRequest.DN.CN, certCertificate.Request.DN.CN)
	assert.Len(t, certCertificate.Request.SAN, 2)
	assert.Equal(t, certCertificate.CACertificate, caCertificateBytes)
	assert.NotEqual(t, certCertificate.Certificate, certCertificateBytes)
	assert.Equal(t, certCertificate.Key, certKeyBytes)

}

func listCertificates(t *testing.T, cli *client.Client) {

	var (
		certificates map[string]client.Certificate
		err          error
	)

	certificates = make(map[string]client.Certificate)

	// 404 - Not found
	_, err = cli.CertificateList("ca-not-found")
	assert.Error(t, err)
	assert.Equal(t, http.StatusText(http.StatusNotFound), err.Error())

	certificates, err = cli.CertificateList(caID)
	assert.Nil(t, err)
	assert.Len(t, certificates, 3)

}

func deleteCertificate(t *testing.T, cli *client.Client) {

	var (
		certificates map[string]client.Certificate

		ok  bool
		err error
	)

	// 404 - Not found
	ok, err = cli.CertificateDelete(caID, "404")
	assert.Error(t, err)
	assert.Equal(t, http.StatusText(http.StatusNotFound), err.Error())
	assert.False(t, ok)

	// 404 - Not found
	ok, err = cli.CertificateDelete("1234", "404")
	assert.Error(t, err)
	assert.Equal(t, http.StatusText(http.StatusNotFound), err.Error())
	assert.False(t, ok)

	// 409 - Conflict - CA certificate cannot be deleted
	ok, err = cli.CertificateDelete(caID, "ca")
	assert.Error(t, err)
	assert.Equal(t, http.StatusText(http.StatusConflict), err.Error())
	assert.False(t, ok)

	// 200 - OK, must match with the certificated created
	ok, err = cli.CertificateDelete(caID, certRequest.DN.CN)
	assert.Nil(t, err)
	assert.True(t, ok)

	// ensure certificate was deleted
	certificates = make(map[string]client.Certificate)

	certificates, err = cli.CertificateList(caID)
	assert.Nil(t, err)
	assert.Len(t, certificates, 2)

}
