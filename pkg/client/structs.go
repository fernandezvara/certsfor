package client

import (
	"crypto/x509"
	"errors"
)

const (
	userAgent = "client.certsfor.dev"
)

// errors
var (
	ErrUnknownError      = errors.New("unknown error connecting to server")
	ErrConnectionRefused = errors.New("connection refused")
)

// Certificate holds the certificate and key file used to interact with the
// data store
type Certificate struct {
	Key             []byte                `json:"key,omitempty"`
	Certificate     []byte                `json:"certificate,omitempty"`
	CACertificate   []byte                `json:"ca_certificate,omitempty"`
	X509Certificate *x509.Certificate     `json:"-"`
	Request         APICertificateRequest `json:"request"`
	CAID            string                `json:"ca_id,omitempty"`
	Parsed          ParsedInfo            `json:"parsed,omitempty"`
}

// APICertificateRequest is the struct with the data needed to create a new
// certificate
type APICertificateRequest struct {
	DN             APIDN    `json:"dn"`
	SAN            []string `json:"san" yaml:"san"`       // SAN
	Key            string   `json:"key" yaml:"key"`       // Key Type (RSA/ECDSA):(complexity)
	ExpirationDays int64    `json:"exp" yaml:"exp"`       // Days the certificate will be valid
	Client         bool     `json:"client" yaml:"client"` // requesting a client certificate?
}

// ParsedInfo holds the information extracted from the certificate if required
type ParsedInfo struct {
	DN             string   `json:"dn"`
	Version        int      `json:"version"`
	SerialNumber   string   `json:"serial_number"`
	NotBefore      int64    `json:"not_before"`
	NotAfter       int64    `json:"not_after"`
	IsCA           bool     `json:"is_ca"`
	DNSNames       []string `json:"dns_names"`
	EmailAddresses []string `json:"emails"`
	IPAddresses    []string `json:"ips"`
	URIs           []string `json:"uris"`
}

// APIDN is the struct of a Distinguished Name
type APIDN struct {
	CN string `json:"cn,omitempty" yaml:"cn"` // common name (required)
	C  string `json:"c,omitempty" yaml:"c"`   // country
	L  string `json:"l,omitempty" yaml:"l"`   // locality
	O  string `json:"o,omitempty" yaml:"o"`   // organization
	OU string `json:"ou,omitempty" yaml:"ou"` // organization unit
	P  string `json:"p,omitempty" yaml:"p"`   // province / state
	PC string `json:"pc,omitempty" yaml:"pc"` // postal code
	ST string `json:"st,omitempty" yaml:"st"` // street
}

// key algorithms
const (
	RSA2048  = "rsa:2048"
	RSA3072  = "rsa:3072"
	RSA4096  = "rsa:4096"
	ECDSA224 = "ecdsa:224"
	ECDSA256 = "ecdsa:256"
	ECDSA384 = "ecdsa:384"
	ECDSA521 = "ecdsa:521"
)

// APIStatus is returned by the API on GET /status
type APIStatus struct {
	Version string `json:"version"`
}
