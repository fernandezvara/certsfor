package firestore

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/fernandezvara/certsfor/db/store"
	"github.com/fernandezvara/rest"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

func init() {
	store.Register("firestore", &Driver{})
}

// Driver initializes a Firestore and returns it
type Driver struct {
}

// Open creates the firestore struct, configures and returns it
func (f Driver) Open(ctx context.Context, connection string) (store.Store, error) {

	var (
		store Firestore
		opt   option.ClientOption
		app   *firebase.App
		err   error
	)

	opt = option.WithCredentialsFile(connection)
	app, err = firebase.NewApp(ctx, &firebase.Config{}, opt)
	if err != nil {
		log.Fatalln(err)
	}

	store.client, err = app.Firestore(ctx)
	if err != nil {
		return nil, err
	}

	return store, err

}

// Firestore is the storage driver to manage data for the Firebase Firestore
//
// collection : CA Identifier that is used to split different datasets on the same storage
//         id : ID of the item to retrieve
//      value : pointer to fill/set with the information
type Firestore struct {
	client *firestore.Client
}

// Get retrieves a value from the storage and unmarshals it to the required type
func (f Firestore) Get(ctx context.Context, collection, id string, value interface{}) (err error) {

	var dsnap *firestore.DocumentSnapshot

	dsnap, err = f.client.Collection(collection).Doc(id).Get(ctx)
	if grpc.Code(err) == codes.NotFound {
		return rest.ErrNotFound
	}

	if err != nil {
		return
	}

	dsnap.DataTo(value)

	return

}

// GetAll retrieves all the items for the required dataset
func (f Firestore) GetAll(ctx context.Context, collection string) (values []map[string]interface{}, err error) {

	var docs []*firestore.DocumentSnapshot

	docs, err = f.client.Collection(collection).Where("ca", "==", false).OrderBy("cn", firestore.Asc).Documents(ctx).GetAll()
	if err != nil {
		return
	}

	for _, doc := range docs {
		values = append(values, doc.Data())
	}

	return
}

// Set inserts/updates a item in the dataset
func (f Firestore) Set(ctx context.Context, collection, id string, value interface{}) (err error) {

	_, err = f.client.Collection(collection).Doc(id).Set(ctx, value)
	return

}

// Delete removes the required ID on the dataset
func (f Firestore) Delete(ctx context.Context, collection, id string) (ok bool, err error) {

	//var result *firestore.WriteResult
	_, err = f.client.Collection(collection).Doc(id).Delete(ctx)

	return
}

// Ping returns a non-nil error if the Store is not healthy or if the
// connection to the persistence is compromised.
func (f Firestore) Ping(ctx context.Context) error { return nil }

// Close releases the resources associated with the Store.
func (f Firestore) Close() error {

	return f.client.Close()

}
