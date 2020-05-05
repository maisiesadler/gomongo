package gomongotesting

import (
	"github.com/digitalfridgedoor/fridgedoorapi/dfdmodels"
	"github.com/digitalfridgedoor/fridgedoordatabase/database"

	"go.mongodb.org/mongo-driver/bson"
)

var overrides = make(map[string]*TestCollection)

// SetTestCollectionOverride sets a the database package to use a TestCollection
func SetTestCollectionOverride() {
	database.SetOverride(overrideDb)
}

// SetTestFindPredicate shows an example of how to override find functionality
func SetTestFindPredicate(predicate func(*dfdmodels.UserView, bson.M) bool) bool {
	fn := func(value interface{}, filter bson.M) bool {
		uv := value.(*dfdmodels.UserView)
		return predicate(uv, filter)
	}

	coll := getOrAddTestCollection("_database", "_collection")
	coll.findPredicate = fn
	return true
}

func overrideDb(database string, collection string) database.ICollection {
	return getOrAddTestCollection(database, collection)
}

func getOrAddTestCollection(database string, collection string) *TestCollection {
	key := database + "_" + collection
	if val, ok := overrides[key]; ok {
		return val
	}
	overrides[key] = CreateTestCollection()
	return overrides[key]
}
