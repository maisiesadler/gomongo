package gomongo

import (
	"context"
	"log"
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// FindByID finds a single result using the collection ID
func (coll *MongoCollection) FindByID(ctx context.Context, objID *primitive.ObjectID, obj interface{}) (interface{}, error) {

	if objID == nil {
		return nil, errNilID
	}

	filter := bson.D{primitive.E{Key: "_id", Value: objID}}

	return coll.FindOne(ctx, filter, obj)
}

// FindOne finds and parses a single result
func (coll *MongoCollection) FindOne(ctx context.Context, filter interface{}, obj interface{}) (interface{}, error) {

	findOneOptions := options.FindOne()

	singleResult := coll.MongoCollection.FindOne(ctx, filter, findOneOptions)

	return parseSingleResult(singleResult, obj)
}

// Find finds all matching results and returns a channel of the parsed result
func (coll *MongoCollection) Find(ctx context.Context, filter interface{}, findOptions *options.FindOptions, obj interface{}) (<-chan interface{}, error) {

	results, err := coll.MongoCollection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}

	r := Parse(ctx, results, obj)
	return r, nil
}

func parseSingleResult(singleResult *mongo.SingleResult, obj interface{}) (interface{}, error) {

	err := singleResult.Err()
	if err != nil {
		return nil, err
	}

	err = singleResult.Decode(obj)

	if err != nil {
		return nil, err
	}

	return obj, nil
}

// Parse parses cursor returned by Find
func Parse(ctx context.Context, cur *mongo.Cursor, obj interface{}) <-chan interface{} {

	objectType := reflect.TypeOf(obj).Elem()

	ch := make(chan interface{})

	go func() {
		defer close(ch)

		// Finding multiple documents returns a cursor
		// Iterating through the cursor allows us to decode documents one at a time
		for cur.Next(ctx) {

			// create a value into which the single document can be decoded
			result := reflect.New(objectType).Interface()
			err := cur.Decode(result)

			if err != nil {
				log.Fatal(err)
			}

			ch <- result
		}

		if err := cur.Err(); err != nil {
			log.Fatal(err)
		}

		// Close the cursor once finished
		cur.Close(ctx)
	}()

	return ch
}
