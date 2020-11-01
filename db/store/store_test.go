package store

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

type storeMock struct{}

func (s storeMock) Get(ctx context.Context, collection, id string, value interface{}) (err error) {
	return nil
}
func (s storeMock) GetAll(ctx context.Context, collection string) (values []map[string]interface{}, err error) {
	return []map[string]interface{}{}, nil
}
func (s storeMock) Set(ctx context.Context, collection, id string, value interface{}) (err error) {
	return nil
}
func (s storeMock) Delete(ctx context.Context, collection, id string) (ok bool, err error) {
	return true, nil
}
func (s storeMock) Ping(ctx context.Context) error {
	return nil
}
func (s storeMock) Close() error {
	return nil
}

// Driver initializes a Store interface and returns it
type driverMock struct{}

func (d driverMock) Open(ctx context.Context, connection string) (Store, error) {

	return storeMock{}, nil

}

func TestInterfaces(t *testing.T) {

	assert.Implements(t, (*Store)(nil), new(storeMock))
	assert.Implements(t, (*Driver)(nil), new(driverMock))

}

func TestOpen(t *testing.T) {

	conn, err := Open(context.Background(), "does-not-exists", "conn")
	assert.Nil(t, conn)
	assert.Error(t, err)

	var driver driverMock

	Register("mock", driver)

	conn, err = Open(context.Background(), "mock", "conn")
	assert.Implements(t, (*Store)(nil), conn)
	assert.Nil(t, err)

	assert.Panics(t, func() { Register("mock", driver) })
	assert.Panics(t, func() { Register("mock1", nil) })

}
