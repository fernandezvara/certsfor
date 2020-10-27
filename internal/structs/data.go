package structs

import "crypto/x509"

// Certificate is the representation of data on the storage
type Certificate struct {
	Key             []byte            `json:"key"`
	Certificate     []byte            `json:"certificate"`
	X509Certificate *x509.Certificate `json:"-"`
}
