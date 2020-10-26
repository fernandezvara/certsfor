package structs

// Certificate is the representation of data on the storage
type Certificate struct {
	Key         []byte `json:"key"`
	Certificate []byte `json:"certificate"`
}
