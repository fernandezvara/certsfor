package client

import "fmt"

// CertificateCreate returns status information if connection is ok to the API, error otherwise
func (c *Client) CertificateCreate(caID, cn string, request APICertificateRequest) (response Certificate, err error) {

	_, err = c.http.Put(fmt.Sprintf("/v1/ca/%s/certificates/%s", caID, cn)).BodyJSON(request).ReceiveSuccess(&response)
	return

}
