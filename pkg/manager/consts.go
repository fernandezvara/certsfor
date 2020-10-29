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
	ErrCommonNameBlank = errors.New("common name cannot be blank")
	ErrCAKeyInvalid    = errors.New("ca key has invalid type")
)
