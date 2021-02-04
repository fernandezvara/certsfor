package client

import (
	"net/http"
)

// Status returns status information if connection is ok to the API, error otherwise
func (c *Client) Status() (status APIStatus, err error) {

	var (
		res *http.Response
	)

	res, err = c.http.Get("/status").ReceiveSuccess(&status)
	err = isError(res, err, http.StatusOK)

	return

}
