package client

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testIsError(t *testing.T) {

	var (
		res *http.Response
		err error
	)

	res = &http.Response{
		StatusCode: 400,
	}

	err = isError(res, http.StatusOK)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), http.StatusText(http.StatusBadRequest))

	err = isError(res, http.StatusOK)
	assert.Nil(t, err)

	// error on *http.Response
	res = &http.Response{
		StatusCode: 400,
	}

	err = isError(res, http.StatusOK)
	assert.Equal(t, err.Error(), http.StatusText(http.StatusBadRequest))

	res = &http.Response{
		StatusCode: 200,
	}

	err = isError(res, http.StatusOK)
	assert.Nil(t, err)

}
