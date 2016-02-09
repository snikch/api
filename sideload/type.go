package sideload

import "reflect"

var typeRegistry = map[reflect.Type]map[int]string{}

// RegisterType registers sideloading information about the type
// of the supplied value for use when retrieving related entities.
func RegisterType(instance interface{}) {
	// Create a lookup map of field index to related entity name.
	fields := map[int]string{}
	// Get the type of the instance supplied.
	typ := reflect.ValueOf(instance).Type()
	for i := 0; i < typ.NumField(); i++ {
		// Loop over every field, and add the sideload entity if it exists.
		field := typ.Field(i)
		name := field.Tag.Get("sideload")
		if name != "" {
			fields[i] = name
		}
	}
	// Register the type.
	typeRegistry[typ] = fields
}
