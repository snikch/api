package sideload

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	// ErrInvalidKind is when the supplied data type is of an invalid kind.
	ErrInvalidKind = errors.New("Invalid type")
)

// Load will load all related entities for the supplied struct, or slice, as
// long as the struct type and callback handlers have been registered.
// Note: A collection is expected to be of the same type, mixed type slices will
// have incorrect data, or may even panic.
func Load(data interface{}, required ...string) (map[string]map[string]interface{}, error) {
	ids, err := idsFromData(data, required...)
	if err != nil {
		return nil, err
	}

	// At this point the ids map is ready for hydration.
	entities, err := hydrateEntitiesFromMap(ids)
	return entities, nil
}

// idsFromData takes a data interface and a list of required fields, and produces
// a lookup map of entity type to entity ids that should be side loaded.
func idsFromData(data interface{}, required ...string) (map[string]map[string]bool, error) {
	value := reflect.ValueOf(data)
	kind := value.Kind()
	switch kind {
	case reflect.Slice:
		fallthrough
	case reflect.Array:
		fallthrough
	case reflect.Struct:
	default:
		// Wrong type - must be array or slice
		return nil, ErrInvalidKind
	}

	ids := map[string]map[string]bool{}
	requiredLookup := map[string]bool{}
	for _, v := range required {
		requiredLookup[v] = true
	}

	// If this is a single entity, just call relatedEntityIDsFromStruct directly.
	if kind == reflect.Struct {
		typ := value.Type()
		fields, ok := typeRegistry[typ]
		// Only load values if this type has been registered.
		if ok {
			relatedEntityIDsFromStruct(value, fields, ids, requiredLookup)
		}
	} else {
		var fields map[int]string
		for i := 0; i <= value.Len()-1; i++ {
			sliceValue := value.Index(i)
			if fields == nil {
				typ := sliceValue.Type()
				var ok bool
				fields, ok = typeRegistry[typ]
				// No need to keep looping, no objects have been registered.
				if !ok {
					break
				}
			}
			relatedEntityIDsFromStruct(sliceValue, fields, ids, requiredLookup)
		}
	}

	return ids, nil
}

// entityCollectionResult represents a single result for a collection of ids.
type entityCollectionResult struct {
	results map[string]interface{}
	name    string
	err     error
}

// relatedEntityIDsFromStruct traverses a single struct instance for fields that
// contain sideload ids. It adds those ids to the ids map if the ids are required.
func relatedEntityIDsFromStruct(val reflect.Value, fields map[int]string, ids map[string]map[string]bool, required map[string]bool) {
	for index, entityName := range fields {
		// If we have required fields, and this entity isn't in it, skip this one.
		if len(required) > 0 && !required[entityName] {
			continue
		}

		// Get the field at the registered index.
		field := val.Field(index)

		// Ensure it's a string type.
		if field.Kind() != reflect.String {
			continue
		}

		// Get the string value.
		id := field.String()
		if id == "" {
			continue
		}
		// Ensure we have a map for this entity type.
		if _, ok := ids[entityName]; !ok {
			ids[entityName] = map[string]bool{}
		}

		// Add the id.
		ids[entityName][id] = true
	}
}

// hydrateEntitiesFromMap will take an entity ids map, and return the results
// for each of the entity types and their ids. This requires a handler for the
// entity types to have been registered previously.
func hydrateEntitiesFromMap(ids map[string]map[string]bool) (map[string]map[string]interface{}, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	entities := map[string]map[string]interface{}{}
	resultChan := make(chan entityCollectionResult)
	for name, idMap := range ids {
		// Turn the map of ids into a slice.
		idSlice := make([]string, len(idMap))
		i := 0
		for id := range idMap {
			idSlice[i] = id
			i++
		}
		// Go get the entity collection
		go singleEntityCollection(resultChan, name, idSlice)
	}

	var firstErr error
	// Loops over the ids again, for the correct count.
	for _ = range ids {
		// Retrieve a result off the result chan.
		result := <-resultChan
		// An error should be assigned to the firstErr is it's empty.
		if result.err != nil {
			if firstErr == nil {
				firstErr = result.err
			}
			continue
		}
		// Results get assigned back to the entities being returned.
		entities[result.name] = result.results
	}
	return entities, firstErr
}

// singleEntityCollection retrieves the results from a single entity handler.
func singleEntityCollection(resultChan chan<- entityCollectionResult, name string, ids []string) {
	// Create a result to send back on the channel
	result := entityCollectionResult{
		name: name,
	}

	// Get the handler for this entity type
	handler, ok := handlerRegistry[name]
	// If we don't have a handler registered, return an error saying so.
	if !ok {
		result.err = fmt.Errorf("No handler registered for entities of type %s", name)
		resultChan <- result
		return
	}

	// Run the handler with the supplied ids, and send it back on the channel.
	result.results, result.err = handler(ids)
	resultChan <- result
}
