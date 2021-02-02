package client

import (
	"fmt"
	"net/http"
)

// CertificateList returns the certificate information if found
func (c *Client) CertificateList(caID string) (response map[string]Certificate, err error) {

	var (
		uri string = fmt.Sprintf("/v1/ca/%s/certificates", caID)
		res *http.Response
	)

	res, err = c.http.Get(uri).ReceiveSuccess(&response)

	err = isError(res, http.StatusOK)

	return

}
