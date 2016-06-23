package changes

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"
)

// SliceDiff represents elements that have been added, removed, or unchanged in
// a slice
type SliceDiff struct {
	Added   []interface{}
	Removed []interface{}
}

type SliceDiffer struct {
	KeyMapper KeyMapper
}

var (
	ErrNotSlice = errors.New("slices must be supplied")
)

func (differ *SliceDiffer) Between(old, new interface{}) (*SliceDiff, error) {
	// No nils thanks.
	if old == nil || new == nil {
		return nil, ErrNil
	}

	// Get the reflect.Values of each, and indirect a pointer if required.
	oldVal := reflect.Indirect(reflect.ValueOf(old))
	newVal := reflect.Indirect(reflect.ValueOf(new))

	// Ensure that each element is a slice.
	if oldVal.Kind() != reflect.Slice || newVal.Kind() != reflect.Slice {
		return nil, ErrNotSlice
	}

	// Loop over each element in old and generate a key lookup
	oldLookup, err := sliceKeys(differ.KeyMapper, oldVal)
	if err != nil {
		return nil, err
	}
	// Loop over each element in new and generate a key lookup
	newLookup, err := sliceKeys(differ.KeyMapper, newVal)
	if err != nil {
		return nil, err
	}

	// Now remove any intersecting keys
	for key := range oldLookup {
		// If this key is also in the new lookup, remove it from both.
		if _, ok := newLookup[key]; ok {
			delete(oldLookup, key)
			delete(newLookup, key)
		}
	}

	// Now generate a slicediff
	diff := &SliceDiff{}
	for _, val := range oldLookup {
		diff.Removed = append(diff.Removed, val)
	}
	for _, val := range newLookup {
		diff.Added = append(diff.Added, val)
	}

	// Retrieve the fields from the mapper.
	return diff, nil
}

func sliceKeys(mapper KeyMapper, val reflect.Value) (map[string]interface{}, error) {
	lookup := map[string]interface{}{}
	for i := 0; i < val.Len(); i++ {
		val := val.Index(i)
		// Get the key indexes and names for changes we care about.
		keyIndexes, err := mapper.KeyIndexes(val)
		if err != nil {
			return nil, err
		}

		// Loop over the key indexes and generate a key for each to match on.
		keyParts := []string{}
		for _, key := range keyIndexes.Keys {
			// Retrieve the field, and then append to the key parts.
			field := val.FieldByIndex(keyIndexes.Indexes[key])
			switch fieldVal := field.Interface().(type) {
			case int, bool, float64, string, time.Time:
				keyParts = append(keyParts, fmt.Sprintf("%v", fieldVal))
				break
			}
		}
		lookup[strings.Join(keyParts, ":")] = val.Interface()
	}
	fmt.Printf("Lookup %+v\n", lookup)
	return lookup, nil
}
