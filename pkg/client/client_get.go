package client

import "fmt"

// CertificateGet returns the certificate information if found
func (c *Client) CertificateGet(caID, cn string, remaining int) (response Certificate, err error) {

	var (
		uri string = fmt.Sprintf("/v1/ca/%s/certificates/%s", caID, cn)
	)

	if remaining > 0 {
		uri = fmt.Sprintf("%s?renew=%d", uri, remaining)
	}

	_, err = c.http.Get(uri).ReceiveSuccess(&response)
	return

}
