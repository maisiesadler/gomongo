package gomongotesting

import (
	"context"
	"errors"

	"github.com/digitalfridgedoor/fridgedoordatabase/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// TestCollection wraps a map
type TestCollection struct {
	coll          map[primitive.ObjectID]interface{}
	findPredicate func(interface{}, bson.M) bool
	setID         func(interface{}, primitive.ObjectID)
}

// CreateTestCollection returns an initialised TestCollection
func CreateTestCollection() *TestCollection {
	coll := make(map[primitive.ObjectID]interface{})
	return &TestCollection{coll: coll}
}

// ensure TestCollection implements interface
func testCollectionIsAnICollection() {
	func(coll database.ICollection) {}(&TestCollection{})
}

func (coll *TestCollection) InsertOne(ctx context.Context, document interface{}) (*primitive.ObjectID, error) {
	id := primitive.NewObjectID()
	coll.coll[id] = document

	// update document with new ID
	if coll.setID != nil {
		coll.setID(document, id)
	}

	return &id, nil
}

func (coll *TestCollection) InsertOneAndFind(ctx context.Context, document interface{}, output interface{}) (interface{}, error) {
	id, err := coll.InsertOne(ctx, document)
	if err != nil {
		return nil, err
	}
	return coll.FindByID(ctx, id, output)
}

func (coll *TestCollection) DeleteByID(ctx context.Context, objID *primitive.ObjectID) error {
	delete(coll.coll, *objID)
	return nil
}

func (coll *TestCollection) UpdateByID(ctx context.Context, objID *primitive.ObjectID, obj interface{}) error {
	coll.coll[*objID] = obj
	return nil
}

func (coll *TestCollection) Find(ctx context.Context, filter interface{}, findOptions *options.FindOptions, obj interface{}) (<-chan interface{}, error) {
	if coll.findPredicate == nil {
		return nil, errors.New("Call SetFindFilter")
	}

	results := make(chan interface{})

	go func() {
		defer close(results)

		elementMap, err := getElementMap(filter)
		if err == nil {
			for _, v := range coll.coll {
				if coll.findPredicate(v, *elementMap) {
					results <- v
				}
			}
		}
	}()

	return results, nil
}

func (coll *TestCollection) FindByID(ctx context.Context, objID *primitive.ObjectID, obj interface{}) (interface{}, error) {
	if o, ok := coll.coll[*objID]; ok {
		return o, nil
	}

	return nil, errors.New("Not found")
}

func (coll *TestCollection) SetFindFilter(predicate func(interface{}, bson.M) bool) {
	coll.findPredicate = predicate
}

func (coll *TestCollection) SetIDSetter(setter func(interface{}, primitive.ObjectID)) {
	coll.setID = setter
}

func (coll *TestCollection) FindOne(ctx context.Context, filter interface{}, obj interface{}) (interface{}, error) {
	if coll.findPredicate == nil {
		return nil, errors.New("Call SetFindFilter")
	}

	elementMap, err := getElementMap(filter)
	if err != nil {
		return nil, err
	}

	for _, v := range coll.coll {
		if coll.findPredicate(v, *elementMap) {
			return v, nil
		}
	}

	return nil, errors.New("No matching element")
}

func getElementMap(filter interface{}) (*bson.M, error) {
	if d, ok := filter.(bson.D); ok {
		m := d.Map()
		return &m, nil
	}
	if m, ok := filter.(bson.M); ok {
		return &m, nil
	}

	return nil, errors.New("filter expected to be bson.D or bson.M")
}
