package client

import (
	"net/http"
)

// CACreate returns status information if connection is ok to the API, error otherwise
func (c *Client) CACreate(request APICertificateRequest) (response Certificate, err error) {

	var (
		res *http.Response
	)

	res, err = c.http.Post("/v1/ca").BodyJSON(request).ReceiveSuccess(&response)
	if err != nil {
		return
	}

	err = isError(res, http.StatusCreated)

	return

}
