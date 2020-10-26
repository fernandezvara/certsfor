package store

import (
	"context"
	"errors"
)

// errors
var (
	ErrDriverNotExists = errors.New("store: driver does not exists")
)

var (
	drivers = make(map[string]Driver)
)

// Register a new store driver for each configured database
func Register(name string, driver Driver) {

	if driver == nil {
		panic("store: Register driver is nil")
	}
	if _, dup := drivers[name]; dup {
		panic("store: Register called twice for driver " + name)
	}
	drivers[name] = driver

}

// Open the connection with the driver, using the connection string with the required data
func Open(ctx context.Context, driver, connection string) (Store, error) {

	if _, dup := drivers[driver]; !dup {
		return nil, ErrDriverNotExists
	}

	return drivers[driver].Open(ctx, connection)

}

// Driver initializes a Store interface and returns it
type Driver interface {
	Open(ctx context.Context, connection string) (Store, error)
}

// Store defines an interface that will be used to interact with the different
// data stores
type Store interface {

	// collection : CA Identifier that is used to split different datasets on the same storage
	//         id : ID of the item to retrieve
	//      value : pointer to fill/set with the information

	// Get retrieves a value from the storage and unmarshals it to the required type
	Get(ctx context.Context, collection, id string, value interface{}) (err error)

	// GetAll retrieves all the items for the required dataset
	GetAll(ctx context.Context, collection string) (values []map[string]interface{}, err error)

	// Set inserts/updates a item in the dataset
	Set(ctx context.Context, collection, id string, value interface{}) (err error)

	// Delete removes the required ID on the dataset
	Delete(ctx context.Context, collection, id string) (ok bool, err error)

	// Ping returns a non-nil error if the Store is not healthy or if the
	// connection to the persistence is compromised.
	Ping(ctx context.Context) error

	// Close releases the resources associated with the Store.
	Close() error
}
