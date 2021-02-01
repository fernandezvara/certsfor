package client

import (
	"fmt"
	"net/http"
)

// CertificateGet returns the certificate information if found
func (c *Client) CertificateGet(caID, cn string, remaining int) (response Certificate, err error) {

	var (
		uri string = fmt.Sprintf("/v1/ca/%s/certificates/%s", caID, cn)
		res *http.Response
	)

	if remaining > 0 {
		uri = fmt.Sprintf("%s?renew=%d", uri, remaining)
	}

	res, err = c.http.Get(uri).ReceiveSuccess(&response)
	if err != nil {
		return
	}

	err = isError(res, http.StatusOK)

	return

}
