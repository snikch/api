package changes

import (
	"reflect"
	"strings"
	"sync"
)

// KeyMapper defines an interface for finding the key name for a given type field.
type KeyMapper interface {
	// KeyIndexes returns the key names, and their location, for comparing. This
	// operation should be heavily cached to avoid runtime performance issues.
	KeyIndexes(reflect.Value) (KeyIndexes, error)
}

// TagMapper is a KeyMapper implements that looks up tags in a sorted order for
// key names, then falls back to the field name.
type TagMapper struct {
	tags  []string
	types map[reflect.Type]KeyIndexes
	sync.RWMutex
}

// KeyIndexes provides an ordered list of keys with their reflection indexes.
type KeyIndexes struct {
	Keys    []string
	Indexes map[string][]int
}

// NewKeyIndexes initialises a new KeyIndexes.
func NewKeyIndexes() KeyIndexes {
	return KeyIndexes{
		Keys:    []string{},
		Indexes: map[string][]int{},
	}
}

// NewTagMapper returns a new TagMapper instance.
func NewTagMapper(tags ...string) *TagMapper {
	return &TagMapper{
		tags:  tags,
		types: map[reflect.Type]KeyIndexes{},
	}
}

// KeyIndexes implements the KeyMapper interface and returns the keys and their
// locations in the value's type.
func (mapper *TagMapper) KeyIndexes(value reflect.Value) (KeyIndexes, error) {
	typ := value.Type()
	if indexes, ok := mapper.types[typ]; ok {
		return indexes, nil
	}
	return mapper.registerValue(value)
}

// registerValue will create an index lookup, save it for later use, and return it.
func (mapper *TagMapper) registerValue(value reflect.Value) (KeyIndexes, error) {
	mapper.Lock()
	indexes := mapper.registerPart("", NewKeyIndexes(), value, []int{})
	mapper.types[value.Type()] = indexes
	mapper.Unlock()
	return indexes, nil
}

func (mapper *TagMapper) registerPart(prefix string, indexes KeyIndexes, val reflect.Value, runningIndex []int) KeyIndexes {
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
		 * and an implementation of a type, we just look at the first level by default.
		 * For example, how do we determine the difference between an embedded struct
		 * and a stdlib struct such as `time.Time`. Chances are we want to compare
		 * the `time.Time` vs the underlying values.
		 * An optional tag can force recursion of a struct. FieldName string `diff:"include"`
		 */
		diffTag := field.Tag.Get("diff")
		if diffTag == "exclude" {
			continue
		}

		// Generate a name for this field, which can be set via a tag.
		var tagName string
		for _, tag := range mapper.tags {
			tagName = field.Tag.Get(tag)
			if tagName != "" {
				break
			}
		}

		var name string
		if tagName != "" {
			name = prefix + tagName
		} else {
			// If recursion of a struct append field name to given prefix.
			name = prefix + field.Name

			// If struct has single field, strip period from prefix and don't append
			// inner field name.
			if typ.NumField() == 1 {
				name = strings.TrimSuffix(prefix, ".")
			}
		}

		// If this is a struct, we can go deeper - but only if it tells us to.
		if fieldType.Kind() == reflect.Struct && diffTag == "include" {
			// IsValid is true if it's not the zero value and CanInterface is true
			// if it's an exported field.
			if !value.IsValid() || !value.CanInterface() {
				continue
			}

			// If this is a struct, recurse down the chain with the interface.
			indexes = mapper.registerPart(name+".", indexes, reflect.ValueOf(value.Interface()), index)
			continue
		}
		indexes.Keys = append(indexes.Keys, name)
		indexes.Indexes[name] = index
	}
	return indexes
}
