package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	_ "github.com/fernandezvara/certsfor/db/badger" // store driver
	"github.com/fernandezvara/certsfor/internal/tests"
	"github.com/fernandezvara/certsfor/pkg/client"
	"github.com/stretchr/testify/assert"
)

var (
	// all comparations
	apiIPPort       string
	caID            string
	caCertificate   []byte
	certCertificate []byte
	certKey         []byte
)

func TestAPI(t *testing.T) {

	apiIPPort = "127.0.0.1:64000"

	// start API in background
	testAPI := tests.TestAPI{}
	go testAPI.StartAPI(t, apiIPPort, []byte{}, []byte{}, []byte{}, false)

	// allow api to start
	time.Sleep(4 * time.Second)

	testStatus(t)            // GET    /status
	testCreateCA(t)          // POST   /v1/ca
	testCreateCertificate(t) // PUT    /v1/ca/:caid/certificates/:cn
	testGetCertificate(t)    // GET    /v1/ca/:caid/certificates/:cn
	testListCertificates(t)  // GET    /v1/ca/:caid/certificates
	testDeleteCertificate(t) // DELETE /v1/ca/:caid/certificates/:cn

	err := testAPI.StopAPI(t)
	assert.Nil(t, err)

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

	assert.Equal(t, tests.Version, response.Version)

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

func testListCertificates(t *testing.T) {

	var (
		response map[string]client.Certificate
		res      *http.Response
		err      error
		certURI  string = uri(fmt.Sprintf("/v1/ca/%s/certificates", caID))
	)

	// 404 - Not found
	res, err = http.Get(uri(fmt.Sprintf("/v1/ca/%s/certificates", "ca-non-existent")))
	assert.Nil(t, err)
	assert.Equal(t, http.StatusNotFound, res.StatusCode)

	// 200 - OK
	res, err = http.Get(certURI)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	err = getFromBody(res, &response)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(response))

}

func testDeleteCertificate(t *testing.T) {

	var (
		status   int
		request  client.APICertificateRequest = requestCertificate(true)
		res      *http.Response
		response map[string]client.Certificate
		err      error
	)

	// 404 - Not found - ca not found
	status, err = sendData(http.MethodDelete, uri(fmt.Sprintf("/v1/ca/%s/certificates/%s", "ca-non-existent", "id-no-existent")), nil, nil)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusNotFound, status)

	// 404 - Not found - id not found
	status, err = sendData(http.MethodDelete, uri(fmt.Sprintf("/v1/ca/%s/certificates/%s", caID, "id-no-existent")), nil, nil)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusNotFound, status)

	// 204 - No content - delete done
	status, err = sendData(http.MethodDelete, uri(fmt.Sprintf("/v1/ca/%s/certificates/%s", caID, request.DN.CN)), nil, nil)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusNoContent, status)

	// 200 - Ok - ensure certificate was deleted
	res, err = http.Get(uri(fmt.Sprintf("/v1/ca/%s/certificates", caID)))
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	err = getFromBody(res, &response)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(response))

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
