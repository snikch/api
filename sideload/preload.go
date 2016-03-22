package sideload

import "github.com/snikch/api/ctx"

type contextKey int

const (
	preloadedEntitiesKey contextKey = iota
)

// Preload adds the supplied entities to the context for addition to the
// sideload payload when it is built. This allows you to manually add additional
// sideloaded entities, and also save additional lookups if you happen to have
// entities available during your application lifecycle.
func Preload(context *ctx.Context, name string, entities map[string]interface{}) {
	// Ignore empty entity maps.
	if len(entities) == 0 {
		return
	}

	context.Lock()
	defer context.Unlock() // Defer in case of panic

	// Get the map of all entities.
	allEntities := preloadedEntities(context)

	if existingEntities, ok := allEntities[name]; ok {
		// If we already have this entity type, add the new ones.
		for id, entity := range entities {
			existingEntities[id] = entity
		}
		allEntities[name] = existingEntities
	} else {
		// If we don't already have this type of entity, just set it.
		allEntities[name] = entities
	}

	// Replace the saved entities with the new set.
	context.SetUnsafe(preloadedEntitiesKey, allEntities)
}

// preloadedEntities returns the entities registered on the context. Note, this
// function performs unsafe operations on the context and should be called with
// a locked context.
func preloadedEntities(context *ctx.Context) map[string]map[string]interface{} {
	if entitiesRaw, ok := context.GetOkUnsafe(preloadedEntitiesKey); ok {
		return entitiesRaw.(map[string]map[string]interface{})
	}

	entities := map[string]map[string]interface{}{}
	context.SetUnsafe(preloadedEntitiesKey, entities)
	return entities
}
