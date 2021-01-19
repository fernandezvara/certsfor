package client

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/dghubble/sling"
)

// Client is the API client if the service is working in client/server mode
type Client struct {
	http *sling.Sling
}

// New returns an API client with default timeout configurations
func New(baseURL, caCert, cert, key string) (*Client, error) {

	return NewWithConnectionTimeouts(baseURL, caCert, cert, key, 5, 5, 10)

}

// NewWithConnectionTimeouts returns a configured API client
func NewWithConnectionTimeouts(baseURL, caCertPath, certPath, keyPath string, dialerTimeout, handshakeTimeout, timeout time.Duration) (*Client, error) {

	var (
		client        Client
		scheme        string = "http"
		httpTransport http.Transport
		certificate   tls.Certificate
		caCertPool    x509.CertPool
		caCertBytes   []byte
		tlsConfig     tls.Config
		err           error
	)

	httpTransport = http.Transport{
		Dial: (&net.Dialer{
			Timeout: dialerTimeout * time.Second,
		}).Dial,
		TLSHandshakeTimeout: handshakeTimeout * time.Second,
	}

	if len(caCertPath) > 0 || len(certPath) > 0 || len(keyPath) > 0 {
		scheme = "https"

		// we can have a mix of configurations so we need to verify both posibilities

		if len(certPath) > 0 && len(keyPath) > 0 {
			certificate, err = tls.LoadX509KeyPair(certPath, keyPath)
			if err != nil {
				return &Client{}, err
			}
			tlsConfig.Certificates = []tls.Certificate{certificate}
		}

		if len(caCertPath) > 0 {
			caCertBytes, err = ioutil.ReadFile(caCertPath)
			if err != nil {
				return &Client{}, err
			}

			caCertPool.AppendCertsFromPEM(caCertBytes)
			tlsConfig.RootCAs = &caCertPool
		}

		httpTransport.TLSClientConfig = &tlsConfig
	}

	client.http = sling.New().Client(&http.Client{
		Timeout:   timeout * time.Second,
		Transport: &httpTransport,
	}).Set("User-Agent", userAgent).Base(fmt.Sprintf("%s://%s/", scheme, baseURL))

	return &client, nil
}
