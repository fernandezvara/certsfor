package client

// CACreate returns status information if connection is ok to the API, error otherwise
func (c *Client) CACreate(request APICertificateRequest) (response Certificate, err error) {

	_, err = c.http.Post("/v1/ca").BodyJSON(request).ReceiveSuccess(&response)
	return

}
