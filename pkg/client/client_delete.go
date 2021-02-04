package client

import (
	"fmt"
	"net/http"
)

// CertificateDelete deletes a certificate
func (c *Client) CertificateDelete(caID, cn string) (ok bool, err error) {

	var (
		uri string = fmt.Sprintf("/v1/ca/%s/certificates/%s", caID, cn)
		res *http.Response
	)

	if res, err = c.http.Delete(uri).ReceiveSuccess(nil); err != nil {
		return
	}

	err = isError(res, err, http.StatusNoContent)
	if err == nil {
		ok = true
	}

	return

}
