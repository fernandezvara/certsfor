package badger

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"runtime"

	"github.com/dgraph-io/badger"
	"github.com/fernandezvara/certsfor/db/store"
	"github.com/fernandezvara/rest"
)

func init() {
	store.Register("badger", &Driver{})
}

// Driver initializes a Badger struct and returns it
type Driver struct {
}

// Open creates the Badger struct, configures and returns it
//
// connection = directory where the database will be stored
func (f Driver) Open(ctx context.Context, connection string) (store.Store, error) {

	var (
		store Badger
		opts  badger.Options
		err   error
	)

	err = os.MkdirAll(connection, 0750)
	if err != nil {
		return nil, err
	}

	opts = badger.DefaultOptions(connection).WithSyncWrites(true).WithLogger(nil)

	if runtime.GOOS == "windows" {
		opts = opts.WithTruncate(true)
	}

	store.db, err = badger.Open(opts)
	if err != nil {
		fmt.Println("err:badger.database.open:", err)
		return nil, err
	}

	return store, err

}

// Badger is the storage driver to manage data using a local database
// (this configuration is not High Available)
type Badger struct {
	db *badger.DB
}

// Get retrieves a value from the storage and unmarshals it to the required type
func (b Badger) Get(ctx context.Context, collection, id string, value interface{}) (err error) {

	var (
		v []byte
	)

	err = b.db.View(func(txn *badger.Txn) error {

		var (
			item *badger.Item
			err  error
		)

		item, err = txn.Get(key(collection, id))
		if err != nil {
			return err
		}

		err = item.Value(func(val []byte) error {
			v = append([]byte{}, val...)
			return nil
		})

		return err

	})

	if err == badger.ErrKeyNotFound {
		return rest.ErrNotFound
	}

	err = json.Unmarshal(v, value)

	return

}

// GetAll retrieves all the items for the required dataset
func (b Badger) GetAll(ctx context.Context, collection string) (values []map[string]interface{}, err error) {

	var fn = func(k, v []byte) error {

		var value map[string]interface{} = make(map[string]interface{})
		err = json.Unmarshal(v, &value)
		if err != nil {
			return err
		}

		values = append(values, value)
		return err

	}

	err = b.db.View(func(txn *badger.Txn) error {

		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 100

		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Seek(key(collection, "")); it.ValidForPrefix(key(collection, "")); it.Next() {

			item := it.Item()
			key := item.Key()
			val, err := item.ValueCopy(nil)
			if err != nil {
				return err
			}

			if err := fn(key, val); err != nil {
				return err
			}

		}
		return nil
	})

	if values == nil {
		err = rest.ErrNotFound
	}

	return

}

// Set inserts/updates a item in the dataset
func (b Badger) Set(ctx context.Context, collection, id string, value interface{}) (err error) {

	var (
		v []byte
	)

	v, err = json.Marshal(value)
	if err != nil {
		return
	}

	return b.db.Update(func(txn *badger.Txn) error {
		e := badger.NewEntry(key(collection, id), v)
		return txn.SetEntry(e)
	})

}

// Delete removes the required ID on the dataset
func (b Badger) Delete(ctx context.Context, collection, id string) (ok bool, err error) {

	err = b.db.Update(func(txn *badger.Txn) error {

		_, err = txn.Get(key(collection, id))
		if err != nil {
			return err
		}

		return txn.Delete(key(collection, id))
	})

	if err == badger.ErrKeyNotFound {
		err = rest.ErrNotFound
	}

	if err == nil {
		ok = true
	}

	return

}

// Ping returns nil (TODO: review)
func (b Badger) Ping(ctx context.Context) error {
	return nil
}

// Close releases the resources associated with the Store.
func (b Badger) Close() error {
	return b.db.Close()
}

func key(collection, id string) []byte {
	return []byte(fmt.Sprintf("%s/%s", collection, id))
}
