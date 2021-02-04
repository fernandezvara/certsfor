package client

import (
	"fmt"
	"net/http"
)

// CertificateCreate returns status information if connection is ok to the API, error otherwise
func (c *Client) CertificateCreate(caID, cn string, request APICertificateRequest) (response Certificate, err error) {

	var (
		res *http.Response
	)

	res, err = c.http.Put(fmt.Sprintf("/v1/ca/%s/certificates/%s", caID, cn)).BodyJSON(request).ReceiveSuccess(&response)
	if err != nil {
		return
	}

	err = isError(res, err, http.StatusOK)

	return

}
