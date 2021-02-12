package client

import (
	"fmt"
	"net/http"
	"net/url"
)

// CertificateGet returns the certificate information if found
func (c *Client) CertificateGet(caID, cn string, remaining int, parse bool) (response Certificate, err error) {

	var (
		uri    string     = fmt.Sprintf("/v1/ca/%s/certificates/%s", caID, cn)
		values url.Values = url.Values{}
		res    *http.Response
	)

	if remaining > 0 {
		values.Set("renew", fmt.Sprintf("%d", remaining))
	}

	if parse {
		values.Set("parse", "true")
	}

	if len(values.Encode()) > 0 {
		uri = fmt.Sprintf("%s?%s", uri, values.Encode())
	}

	res, err = c.http.Get(uri).ReceiveSuccess(&response)
	if err != nil {
		return
	}

	err = isError(res, err, http.StatusOK)

	return

}
