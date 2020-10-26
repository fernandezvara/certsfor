package structs

// Config is the struct that contains the configuration for the service
type Config struct {
	Store StoreConfig `json:"store" yaml:"store"`
}

// StoreConfig is the struct that contains the configuration for the store and its connection string
type StoreConfig struct {
	Type             string `json:"type" yaml:"type"`
	ConnectionString string `json:"connection_string" yaml:"connection_string"`
}
