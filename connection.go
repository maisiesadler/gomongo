package gomongo

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClient *mongo.Client

func connected() bool {
	return mongoClient != nil
}

func connect(ctx context.Context, connectionString string) bool {
	// Set client options
	clientOptions := options.Client().ApplyURI(connectionString)

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(ctx, nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	mongoClient = client

	return true
}

// Disconnect removes the current connection to mongo
func Disconnect() error {
	if mongoClient == nil {
		return nil
	}
	err := mongoClient.Disconnect(context.TODO())
	mongoClient = nil

	return err
}
