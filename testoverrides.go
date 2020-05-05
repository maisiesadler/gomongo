package gomongo

var createOverrideFn func(string, string) ICollection

func tryGetOverrideFor(database string, collection string) (ICollection, bool) {

	if createOverrideFn != nil {
		override := createOverrideFn(database, collection)
		return override, true
	}

	return nil, false
}

// SetOverride allows an override ICollection to be set for testing
func SetOverride(overrideFn func(string, string) ICollection) {
	createOverrideFn = overrideFn
}
