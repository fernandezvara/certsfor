package tests

import (
	"io/ioutil"
	"testing"
	"time"

	_ "github.com/fernandezvara/certsfor/db/badger" // store driver
	"github.com/stretchr/testify/assert"
)

func TestAPIForTesting(t *testing.T) {

	const localIPPort = "127.0.0.1:63999"

	testAPI := TestAPI{}

	testAPI.StartAPI(t, localIPPort, []byte{}, []byte{}, []byte{})

	time.Sleep(5 * time.Second)

	err := testAPI.StopAPI(t)
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
