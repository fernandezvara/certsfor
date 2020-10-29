package structs

import "crypto/x509"

// Certificate holds the certificate and key file used to interact with the
// data store
type Certificate struct {
	Key             []byte            `json:"key"`
	Certificate     []byte            `json:"certificate"`
	X509Certificate *x509.Certificate `json:"-"`
}

// APICertificateRequest is the struct with the data needed to create a new
// certificate
type APICertificateRequest struct {
	CN             string   `json:"cn,omitempty"` // common name (required)
	C              string   `json:"c,omitempty"`  // country
	L              string   `json:"l,omitempty"`  // locality
	O              string   `json:"o,omitempty"`  // organization
	OU             string   `json:"ou,omitempty"` // organization unit
	P              string   `json:"p,omitempty"`  // province
	PC             string   `json:"pc,omitempty"` // postal code
	ST             string   `json:"st,omitempty"` // street
	SAN            []string `json:"san"`          // SAN
	Key            string   `json:"key"`          // Key Type (RSA/ECDSA):(complexity)
	ExpirationDays int64    `json:"expiration"`   // Days the certificate will be valid
}

// key algorithms
const (
	RSA2048  = "rsa:2048"
	RSA4096  = "rsa:4096"
	ECDSA224 = "ecdsa:224"
	ECDSA256 = "ecdsa:256"
	ECDSA384 = "ecdsa:384"
	ECDSA521 = "ecdsa:521"
)
