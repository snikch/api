package changes

import "reflect"

// KeyMapper defines an interface for finding the key name for a given type field.
type KeyMapper interface {
	// KeyIndexes returns the key names, and their location, for comparing. This
	// operation should be heavily cached to avoid runtime performance issues.
	KeyIndexes(reflect.Value) (map[string][]int, error)
}

// TagMapper is a KeyMapper implements that looks up tags in a sorted order for
// key names, then falls back to the field name.
type TagMapper struct {
	tags  []string
	types map[reflect.Type]map[string][]int
}

// NewTagMapper returns a new TagMapper instance.
func NewTagMapper(tags ...string) *TagMapper {
	return &TagMapper{
		tags:  tags,
		types: map[reflect.Type]map[string][]int{},
	}
}

// KeyIndexes implements the KeyMapper interface and returns the keys and their
// locations in the value's type.
func (mapper *TagMapper) KeyIndexes(value reflect.Value) (map[string][]int, error) {
	typ := value.Type()
	if indexes, ok := mapper.types[typ]; ok {
		return indexes, nil
	}
	return mapper.registerValue(value)
}

// registerValue will create an index lookup, save it for later use, and return it.
func (mapper *TagMapper) registerValue(value reflect.Value) (map[string][]int, error) {
	indexes := mapper.registerPart(map[string][]int{}, value, []int{})
	mapper.types[value.Type()] = indexes
	return indexes, nil
}

func (mapper *TagMapper) registerPart(indexes map[string][]int, val reflect.Value, runningIndex []int) map[string][]int {
	// Get the type of the instance supplied.
	typ := val.Type()
	for i := 0; i < typ.NumField(); i++ {
		// Loop over every field, and add the sideload entity if it exists.
		field := typ.Field(i)
		// Exclude unexported fields.
		if field.PkgPath != "" {
			continue
		}

		fieldType := field.Type
		value := val.Field(i)
		// If this is a pointer, we need to take the pointer field and value.
		if fieldType.Kind() == reflect.Ptr {
			fieldType = fieldType.Elem()
			value = value.Elem()
		}

		// Create the index value.
		index := make([]int, len(runningIndex)+1)
		copy(index, runningIndex)
		index[len(runningIndex)] = i
		/**
		 * Due to the complicity of differentiating between a straight up struct,
		 * and an implementation of a type, we just look at the first level.
		 * For example, how do we determine the difference between an embedded struct
		 * and a stdlib struct such as `time.Time`. Chances are we want to compare
		 * the `time.Time` vs the underlying values.
		 */

		// // If this is a struct, we can go deeper.
		// if fieldType.Kind() == reflect.Struct {
		// 	// IsValid is true if it's not the zero value and CanInterface is true
		// 	// if it's an exported field.
		// 	if !value.IsValid() || !value.CanInterface() {
		// 		continue
		// 	}
		// 	// If this is a struct, recurse down the chain with the interface.
		// 	indexes = mapper.registerPart(indexes, reflect.ValueOf(value.Interface()), index)
		// 	continue
		// }
		// switch i := reflect.Indirect(value).Interface().(type) {
		// case string, *string, time.Time, *time.Time, int, *int, float64, *float64, bool, *bool:
		// Determine the name by trying tags, then just using the field name.
		// fmt.Println("Worked on", i)
		var name string
		for _, tag := range mapper.tags {
			name = field.Tag.Get(tag)
			if name != "" {
				break
			}
		}
		if name == "" {
			name = field.Name
		}
		indexes[name] = index
		// default:
		// 	fmt.Println("Failed on", i)
		// }
	}
	return indexes
}
