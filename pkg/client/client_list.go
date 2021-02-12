package client

import (
	"fmt"
	"net/http"
)

// CertificateList returns the certificate information if found
func (c *Client) CertificateList(caID string, parseInfo bool) (response map[string]Certificate, err error) {

	var (
		uri string = fmt.Sprintf("/v1/ca/%s/certificates", caID)
		res *http.Response
	)

	if parseInfo {
		uri = fmt.Sprintf("%s?parse=true", uri)
	}

	res, err = c.http.Get(uri).ReceiveSuccess(&response)
	err = isError(res, err, http.StatusOK)

	return

}
