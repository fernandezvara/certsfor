package manager

import "errors"

// File literals constants
const (
	FileRSAPrivateKey = "RSA PRIVATE KEY"
	FileCertificate   = "CERTIFICATE"
)

// Errors
var (
	ErrUnparseableFile = errors.New("unparseable file")
)
