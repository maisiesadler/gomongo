package gomongo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// InsertOne inserts a new document into the collection
func (coll *MongoCollection) InsertOne(ctx context.Context, document interface{}) (*primitive.ObjectID, error) {
	insertOneOptions := options.InsertOne()

	insertOneResult, err := coll.MongoCollection.InsertOne(ctx, document, insertOneOptions)
	if err != nil {
		return nil, err
	}

	insertedID := insertOneResult.InsertedID.(primitive.ObjectID)
	return &insertedID, nil
}

// InsertOneAndFind inserts a new document into the collection and returns the new object
func (coll *MongoCollection) InsertOneAndFind(ctx context.Context, document interface{}, output interface{}) (interface{}, error) {

	insertedID, err := coll.InsertOne(ctx, document)
	if err != nil {
		return nil, err
	}

	return coll.FindByID(ctx, insertedID, output)
}

// DeleteByID deletes a document by id
func (coll *MongoCollection) DeleteByID(ctx context.Context, objID *primitive.ObjectID) error {

	deleteOptions := options.Delete()

	_, err := coll.MongoCollection.DeleteOne(ctx, bson.D{primitive.E{Key: "_id", Value: objID}}, deleteOptions)
	if err != nil {
		return err
	}

	return nil
}

// UpdateByID finds and updates an object by ID
func (coll *MongoCollection) UpdateByID(ctx context.Context, objID *primitive.ObjectID, obj interface{}) error {

	o := options.FindOneAndReplace()

	filter := bson.D{primitive.E{Key: "_id", Value: objID}}

	singleResult := coll.MongoCollection.FindOneAndReplace(ctx, filter, obj, o)

	return singleResult.Err()
}

// UpdateByIDAndFind finds and updates an object by ID and returns the new object
func (coll *MongoCollection) UpdateByIDAndFind(ctx context.Context, objID *primitive.ObjectID, obj interface{}) (interface{}, error) {

	err := coll.UpdateByID(ctx, objID, obj)
	if err != nil {
		return nil, err
	}

	return coll.FindByID(ctx, objID, obj)
}
