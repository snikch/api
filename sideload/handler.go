package sideload

import "github.com/snikch/api/ctx"

var handlerRegistry = map[string]EntityHandler{}

// EntityHandler represents a callback handler for an entity type.
type EntityHandler func(*ctx.Context, []string) (map[string]interface{}, error)

// RegisterEntityHandler registers an EntityHandler for handling a specific type
// of entity.
func RegisterEntityHandler(name string, handler EntityHandler) {
	handlerRegistry[name] = handler
}
