package tests

import (
	"context"
	"io/ioutil"
	"testing"
	"time"

	_ "github.com/fernandezvara/certsfor/db/badger" // store driver
	"github.com/fernandezvara/certsfor/db/store"
	"github.com/fernandezvara/certsfor/internal/service"
	"github.com/stretchr/testify/assert"
)

func TestAPIForTesting(t *testing.T) {

	const localIPPort = "127.0.0.1:63999"

	testAPI := TestAPI{}

	testAPI.StartAPI(t, localIPPort, []byte{}, []byte{}, []byte{}, false)

	time.Sleep(1 * time.Second)

	err := testAPI.StopAPI(t)
	assert.Nil(t, err)

}

func TestAPIWithService(t *testing.T) {

	const localIPPort = "127.0.0.1:63998"

	var (
		databaseDir string
		sto         store.Store
		srv         *service.Service
		err         error
	)

	// create temporal directory for the database
	databaseDir, err = ioutil.TempDir("", "cfd")
	assert.Nil(t, err)

	sto, err = store.Open(context.Background(), "badger", databaseDir)
	assert.Nil(t, err)

	srv = service.NewAsServer(sto, Version)

	testAPI := TestAPI{}

	testAPI.StartAPIWithService(t, localIPPort, []byte{}, []byte{}, []byte{}, false, srv)

	time.Sleep(1 * time.Second)

	err = testAPI.StopAPI(t)
	assert.Nil(t, err)

}

func TestCleanup(t *testing.T) {

	var (
		databaseDir string
		err         error
	)

	// create temporal directory for the database
	databaseDir, err = ioutil.TempDir("", "cfd")
	assert.Nil(t, err)

	assert.DirExists(t, databaseDir)
	cleanup(databaseDir)

	assert.NoDirExists(t, databaseDir)

}
