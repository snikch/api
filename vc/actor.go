package vc

import "fmt"

// Actor represents a request actor, and is generally either an oauth client or user.
type Actor interface {
	// ActorInfo returns both the identifier and type of the actor.
	ActorInfo() (actorID string, actorType string)
}

// ActorID returns an identifier for the supplied actor.
func ActorID(actor Actor) string {
	id, typ := actor.ActorInfo()
	return fmt.Sprintf("%s:%s", id, typ)
}
