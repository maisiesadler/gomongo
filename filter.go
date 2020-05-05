package gomongo

import "go.mongodb.org/mongo-driver/bson/primitive"

// Filter returns an array of objectIDs matching the filterFn
func Filter(ids []primitive.ObjectID, filterFn func(id *primitive.ObjectID) bool) []primitive.ObjectID {
	filtered := []primitive.ObjectID{}

	for id := range iterateObjectIDs(ids) {
		if filterFn(id) {
			filtered = append(filtered, *id)
		}
	}

	return filtered
}

func iterateObjectIDs(ids []primitive.ObjectID) <-chan *primitive.ObjectID {
	ch := make(chan *primitive.ObjectID)

	go func() {
		defer close(ch)
		for _, id := range ids {
			ch <- &id
		}
	}()

	return ch
}
