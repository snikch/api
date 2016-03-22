package sideload

import "reflect"

var (
	typeRegistry = map[reflect.Type]map[string][][]int{}
)

// RegisterType registers sideloading information about the type
// of the supplied value for use when retrieving related entities.
func RegisterType(instance interface{}) {
	// Create a lookup map of field index to related entity name.
	fields := map[string][][]int{}
	index := []int{}

	// Register the type.
	typ := reflect.ValueOf(instance).Type()
	typeRegistry[typ] = registerPart(fields, instance, index)
}

func registerPart(fields map[string][][]int, instance interface{}, index []int) map[string][][]int {
	// Get the type of the instance supplied.
	val := reflect.ValueOf(instance)
	typ := val.Type()
	for i := 0; i < typ.NumField(); i++ {
		// Loop over every field, and add the sideload entity if it exists.
		field := typ.Field(i).Type
		value := val.Field(i)
		// If this is a pointer, we need to take the pointer field and value.
		if field.Kind() == reflect.Ptr {
			field = field.Elem()
			value = value.Elem()
		}

		switch field.Kind() {
		case reflect.Struct:
			// IsValid is true if it's not the zero value and CanInterface is true
			// if it's an exported field.
			if !value.IsValid() || !value.CanInterface() {
				continue
			}
			lookup := make([]int, len(index)+1)
			copy(lookup, index)
			lookup[len(index)] = i
			// If this is a struct, recurse down the chain with the interface.
			fields = registerPart(fields, value.Interface(), lookup)
		case reflect.String:
			// If this is a string with a sideload value, we can just take the
			// index address for it.
			name := typ.Field(i).Tag.Get("sideload")
			if name != "" {
				field, ok := fields[name]
				if !ok {
					field = [][]int{}
				}
				lookup := make([]int, len(index)+1)
				copy(lookup, index)
				lookup[len(index)] = i
				field = append(field, lookup)
				fields[name] = field
			}
		}
	}
	return fields
}
