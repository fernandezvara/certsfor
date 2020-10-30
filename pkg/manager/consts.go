package manager

import "errors"

// File literals constants
const (
	FilePrivateKey  = "PRIVATE KEY"
	FileCertificate = "CERTIFICATE"
)

// Errors
var (
	ErrUnparseableFile = errors.New("unparseable file")
	ErrCommonNameBlank = errors.New("common name cannot be blank")
	ErrKeyInvalid      = errors.New("key has invalid type")
)
