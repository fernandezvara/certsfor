package tests

import (
	"context"
	"io/ioutil"
	"os"
	"testing"

	"github.com/fernandezvara/certsfor/db/store"
	"github.com/fernandezvara/certsfor/internal/api"
	"github.com/fernandezvara/certsfor/internal/service"
	"github.com/stretchr/testify/assert"
)

const (
	// Version used on test server
	Version = "test-version"
)

// TestAPI is a simple struct to start and stop an API for testing
type TestAPI struct {
	api         *api.API
	databaseDir string
}

// StartAPI is exported to be useful on other libraries tests. Starts the api in testing mode
func (ta *TestAPI) StartAPI(t *testing.T, apiIPPort string, cert, key, caCert []byte, requireClientCert bool) {

	var (
		sto store.Store
		srv *service.Service
		err error
	)

	// create temporal directory for the database
	ta.databaseDir, err = ioutil.TempDir("", "cfd")
	assert.Nil(t, err)

	sto, err = store.Open(context.Background(), "badger", ta.databaseDir)
	assert.Nil(t, err)

	srv = service.NewAsServer(sto, Version)
	ta.api = api.New(srv, Version)
	go ta.api.Start(apiIPPort, cert, key, caCert, requireClientCert, []string{"stdout"}, []string{"stdout"}, true)

}

// StartAPIWithService is exported to be useful on other libraries tests. Starts the api in testing mode
func (ta *TestAPI) StartAPIWithService(t *testing.T, apiIPPort string, cert, key, caCert []byte, requireClientCert bool, srv *service.Service) {

	ta.api = api.New(srv, Version)
	go ta.api.Start(apiIPPort, cert, key, caCert, requireClientCert, []string{"stdout"}, []string{"stdout"}, true)

}

// StopAPI stops the API for testing
func (ta *TestAPI) StopAPI(t *testing.T) (err error) {

	defer cleanup(ta.databaseDir)

	err = ta.api.Stop()
	assert.Nil(t, err)

	return

}

func cleanup(databaseDir string) {

	os.RemoveAll(databaseDir)

}
