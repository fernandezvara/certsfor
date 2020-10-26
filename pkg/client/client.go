package client

import "github.com/dghubble/sling"

// Client is the API client if the service is working in client/server mode
type Client struct {
	http *sling.Sling
}

// New returns an API client
func New() (*Client, error) {
	return &Client{}, nil
}

// Ping returns nil if connection is ok to the API, error otherwise
func (c *Client) Ping() error {

	return nil

}
