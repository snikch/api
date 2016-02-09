package sideload

var handlerRegistry = map[string]EntityHandler{}

// EntityHandler represents a callback handler for an entity type.
type EntityHandler func([]string) (map[string]interface{}, error)

// RegisterEntityHandler registers an EntityHandler for handling a specific type
// of entity.
func RegisterEntityHandler(name string, handler EntityHandler) {
	handlerRegistry[name] = handler
}
