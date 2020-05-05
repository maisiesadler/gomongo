package gomongotesting

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TestType struct {
	ID       primitive.ObjectID
	Username string
}

func TestCanFindOne(t *testing.T) {

	// Arrange
	username := "Test"
	test := &TestType{
		Username: username,
	}

	coll := CreateTestCollection()
	coll.SetFindFilter(example)

	_, err := coll.InsertOne(context.TODO(), test)
	assert.Nil(t, err, "Error inserting to collection")

	// Act
	found, err := coll.FindOne(context.TODO(), bson.D{primitive.E{Key: "username", Value: username}}, &TestType{})

	// Assert
	assert.Nil(t, err, "Error finding document")
	assert.NotNil(t, found, "Did not find document")
	assert.Equal(t, username, found.(*TestType).Username)
}

func TestCanFindMultiple(t *testing.T) {

	// Arrange
	username := "Test"
	test := &TestType{
		Username: username,
	}

	coll := CreateTestCollection()
	coll.SetFindFilter(example)

	_, err := coll.InsertOne(context.TODO(), test)
	assert.Nil(t, err, "Error inserting to collection")

	coll.InsertOne(context.TODO(), &TestType{
		Username: username,
	})
	coll.InsertOne(context.TODO(), &TestType{
		Username: "another name",
	})

	findOptions := options.Find()
	filter := bson.D{primitive.E{Key: "username", Value: username}}

	// Assert
	results, err := coll.Find(context.TODO(), filter, findOptions, &TestType{})
	assert.Nil(t, err, "Error finding document")
	assert.NotNil(t, results, "Results nil")

	parsed := make([]*TestType, 0)

	for i := range results {
		parsed = append(parsed, i.(*TestType))
	}

	if len(parsed) != 2 {
		t.Errorf("Did not get expected number of results. Expected 2, Actual %v.", len(parsed))
	}
}

func TestCanAssignObjectID(t *testing.T) {

	// Arrange
	ctx := context.Background()
	username := "Test"
	test := &TestType{
		Username: username,
	}

	coll := CreateTestCollection()
	coll.SetIDSetter(func(document interface{}, id primitive.ObjectID) {
		// try update id on model
		if u, ok := document.(*TestType); ok {
			u.ID = id
		}
	})

	insertedID, err := coll.InsertOne(ctx, test)
	updated, finderr := coll.FindByID(ctx, insertedID, &TestType{})

	// Assert
	assert.Nil(t, err, "Error inserting to collection")
	assert.Nil(t, finderr, "Error finding updated result")
	assert.NotNil(t, insertedID)
	assert.Equal(t, *insertedID, test.ID)
	updatedTestType, ok := updated.(*TestType)
	assert.True(t, ok)
	assert.Equal(t, *insertedID, updatedTestType.ID)
}

func example(value interface{}, filter bson.M) bool {
	test := value.(*TestType)
	return filter["username"] == test.Username
}
