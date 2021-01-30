package api_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"

	_ "github.com/fernandezvara/certsfor/db/badger" // store driver
	"github.com/fernandezvara/certsfor/db/store"
	"github.com/fernandezvara/certsfor/internal/api"
	"github.com/fernandezvara/certsfor/internal/service"
	"github.com/fernandezvara/certsfor/pkg/client"
	"github.com/stretchr/testify/assert"
)

const (
	version   = "test-version"
	apiIPPort = "127.0.0.1:64000"
)

var (
	// all comparations
	caID            string
	caCertificate   []byte
	certCertificate []byte
	certKey         []byte
)

func TestAPI(t *testing.T) {

	// start API in background
	go startAPI(t)

	// allow api to start
	time.Sleep(2 * time.Second)

	testStatus(t)            // GET  /status
	testCreateCA(t)          // POST /v1/ca
	testCreateCertificate(t) // PUT  /v1/ca/:caid/certificates/:cn
	testGetCertificate(t)    // GET  /v1/ca/:caid/certificates/:cn

	killAPI(t)

}

func testStatus(t *testing.T) {

	var (
		res      *http.Response
		response client.APIStatus
		err      error
	)

	res, err = http.Get(uri("/status"))
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	err = getFromBody(res, &response)
	assert.Nil(t, err)

	assert.Equal(t, version, response.Version)

}

func testCreateCA(t *testing.T) {

	var (
		request  client.APICertificateRequest
		response client.Certificate
		status   int
		err      error
	)

	// 405 - Method Not Allowed
	status, err = sendData(http.MethodPut, uri("/v1/ca"), request, nil)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusMethodNotAllowed, status)

	// 400 - Bad Request (request without required input must fail)
	status, err = sendData(http.MethodPost, uri("/v1/ca"), request, nil)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, status)

	// 400 - Bad Request (body malformed)
	status, err = sendData(http.MethodPost, uri("/v1/ca"), []byte(`{ "this": "is malformed}`), nil)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, status)

	request.DN.CN = "test-ca"
	request.SAN = []string{"ca.test.example.com", "192.168.1.1"}
	request.Key = client.ECDSA521
	request.ExpirationDays = 365

	// 200 - OK
	status, err = sendData(http.MethodPost, uri("/v1/ca"), request, &response)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, status)
	assert.Equal(t, request.DN.CN, response.Request.DN.CN)
	assert.Equal(t, request.SAN, response.Request.SAN)
	assert.Equal(t, request.Key, response.Request.Key)
	assert.Equal(t, request.ExpirationDays, response.Request.ExpirationDays)
	assert.Greater(t, len(response.Certificate), 800)
	assert.Greater(t, len(response.Key), 300)
	assert.NotEmpty(t, response.CAID)
	caID = response.CAID
	caCertificate = response.Certificate

}

func requestCertificate(ok bool) (request client.APICertificateRequest) {

	request.DN.CN = "cert1"
	if ok {
		request.DN.CN = "cert"
	}
	request.DN.O = "org"
	request.DN.OU = "ou"
	request.SAN = []string{"service1.test.example.com", "192.168.1.2"}
	request.Key = client.ECDSA521
	request.ExpirationDays = 90

	return

}

func testCreateCertificate(t *testing.T) {

	var (
		request  client.APICertificateRequest
		response client.Certificate
		status   int
		err      error
		certURI  string = uri(fmt.Sprintf("/v1/ca/%s/certificates/cert", caID))
	)

	// 405 - Method Not Allowed
	status, err = sendData(http.MethodPost, certURI, request, nil)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusMethodNotAllowed, status)

	// 400 - Bad Request (request without required input must fail)
	status, err = sendData(http.MethodPut, certURI, request, nil)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, status)

	// 400 - Bad Request (body malformed)
	status, err = sendData(http.MethodPut, certURI, []byte(`{ "this": "is malformed}`), nil)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, status)

	// 409 - Conflict (ca certificate cannot be overwritten)
	status, err = sendData(http.MethodPut, uri(fmt.Sprintf("/v1/ca/%s/certificates/ca", caID)), request, nil)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusConflict, status)

	request = requestCertificate(false)

	// 409 - Conflict (cn on url != cert.DN.CN)
	status, err = sendData(http.MethodPut, certURI, request, nil)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusConflict, status)

	request = requestCertificate(true)

	// 200 - OK
	status, err = sendData(http.MethodPut, certURI, request, &response)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, status)
	assert.Equal(t, request.DN.CN, response.Request.DN.CN)
	assert.Equal(t, request.DN.O, response.Request.DN.O)
	assert.Equal(t, request.DN.OU, response.Request.DN.OU)
	assert.Equal(t, request.SAN, response.Request.SAN)
	assert.Equal(t, request.Key, response.Request.Key)
	assert.Equal(t, request.ExpirationDays, response.Request.ExpirationDays)
	assert.Greater(t, len(response.Certificate), 800)
	assert.Greater(t, len(response.Key), 300)
	assert.Equal(t, caCertificate, response.CACertificate)

	certCertificate, certKey = response.Certificate, response.Key

}

func testGetCertificate(t *testing.T) {

	var (
		request  client.APICertificateRequest
		response client.Certificate
		res      *http.Response
		err      error
		certURI  string = uri(fmt.Sprintf("/v1/ca/%s/certificates/cert", caID))
	)

	request = requestCertificate(true)

	// 404 - Not found
	res, err = http.Get(uri(fmt.Sprintf("/v1/ca/%s/certificates/404", caID)))
	assert.Nil(t, err)
	assert.Equal(t, http.StatusNotFound, res.StatusCode)

	// 400 - Bad Request (renew is not a number)
	res, err = http.Get(uri(fmt.Sprintf("/v1/ca/%s/certificates/cert?renew=NOT_A_NUMBER", caID)))
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)

	// 200 - OK
	res, err = http.Get(certURI)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	err = getFromBody(res, &response)
	assert.Nil(t, err)
	assert.Equal(t, request.DN.CN, response.Request.DN.CN)
	assert.Equal(t, request.DN.O, response.Request.DN.O)
	assert.Equal(t, request.DN.OU, response.Request.DN.OU)
	assert.Equal(t, request.SAN, response.Request.SAN)
	assert.Equal(t, request.Key, response.Request.Key)
	assert.Equal(t, request.ExpirationDays, response.Request.ExpirationDays)
	assert.Greater(t, len(response.Certificate), 800)
	assert.Greater(t, len(response.Key), 300)
	assert.Equal(t, caCertificate, response.CACertificate)
	assert.Equal(t, certCertificate, response.Certificate)
	assert.Equal(t, certKey, response.Key)

	// 200 - OK + Force Renew
	res, err = http.Get(uri(fmt.Sprintf("/v1/ca/%s/certificates/cert?renew=100", caID)))
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	err = getFromBody(res, &response)
	assert.Nil(t, err)
	assert.Equal(t, request.DN.CN, response.Request.DN.CN)
	assert.Equal(t, request.DN.O, response.Request.DN.O)
	assert.Equal(t, request.DN.OU, response.Request.DN.OU)
	assert.Equal(t, request.SAN, response.Request.SAN)
	assert.Equal(t, request.Key, response.Request.Key)
	assert.Equal(t, request.ExpirationDays, response.Request.ExpirationDays)
	assert.Equal(t, caCertificate, response.CACertificate)
	assert.NotEqual(t, certCertificate, response.Certificate)
	assert.Equal(t, certKey, response.Key)

}

func startAPI(t *testing.T) {

	var (
		databaseDir string
		sto         store.Store
		srv         *service.Service
		testAPI     *api.API
		err         error
	)

	// create temporal directory for the database
	databaseDir, err = ioutil.TempDir("", "cdf")
	fmt.Println(databaseDir, err)
	assert.Nil(t, err)
	defer cleanup(databaseDir)

	sto, err = store.Open(context.Background(), "badger", databaseDir)
	assert.Nil(t, err)

	srv = service.NewAsServer(sto, version)
	testAPI = api.New(srv, version)
	err = testAPI.Start(apiIPPort, []byte{}, []byte{}, []byte{}, []string{"stdout"}, []string{"stdout"}, true)
	assert.Nil(t, err)

}

func killAPI(t *testing.T) {

	process, err := os.FindProcess(os.Getegid())
	assert.Nil(t, err)

	process.Signal(os.Interrupt)

}

func cleanup(databaseDir string) {

	err := os.RemoveAll(databaseDir)
	if err != nil {
		panic(err)
	}

}

// helpers
func uri(path string) string {
	return fmt.Sprintf("http://%s%s", apiIPPort, path)
}

func getFromBody(res *http.Response, obj interface{}) error {

	defer res.Body.Close()
	objectByte, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(objectByte, &obj)

}

func sendData(method, uri string, request interface{}, response interface{}) (status int, err error) {

	var (
		client    *http.Client
		req       *http.Request
		res       *http.Response
		bodyBytes []byte
	)

	client = http.DefaultClient

	switch request.(type) {
	case []byte:
		bodyBytes = request.([]byte)
	default:
		bodyBytes, err = json.Marshal(request)
		if err != nil {
			return
		}
	}

	req, err = http.NewRequest(method, uri, bytes.NewReader(bodyBytes))
	if err != nil {
		return
	}

	req.Header.Add("Content-Type", "application/json")
	res, err = client.Do(req)
	if err != nil {
		return
	}

	status = res.StatusCode
	if response != nil {
		err = getFromBody(res, response)
	}

	return

}
