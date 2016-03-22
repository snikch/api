package vc

import (
	"net/http"
	"time"

	"github.com/snikch/api/ctx"
)

// Criteria defines how results should be modified, or filtered.
type Criteria struct {
	// Sideload is a slice of related entities to sideload.
	Sideload []string
	// OrderColumn defines the order that collections should be ordered by.
	OrderColumn    *string
	OrderAscending bool
	// From represents a time that entities should have been updated after.
	From *time.Time
	// Limit is the number of results in a collection
	Limit int
	// State is an array of item states to filter by
	State []string
}

// SideloadQueryKey is the name of the query parameter that determines which
// entities should be sideloaded.
var SideloadQueryKey = "include"

// StateQueryKey is the name of the query paramtere that determines which
// entity states should be returned.
var StateQueryKey = "state"

// RequestCriteria generates a criteria instance from the request.
func RequestCriteria(r *http.Request) *Criteria {
	criteria := Criteria{
		Sideload: r.URL.Query()[SideloadQueryKey],
		State:    r.URL.Query()[StateQueryKey],
	}
	// Prevent overloading of the sideload and state criteria.
	if len(criteria.Sideload) > 100 {
		criteria.Sideload = nil
	}
	if len(criteria.State) > 100 {
		criteria.State = nil
	}
	return &criteria
}

// SetContextCriteria sets the criteria against a context.
func SetContextCriteria(context *ctx.Context, criteria *Criteria) {
	context.Set(criteriaContextKey, criteria)
}

// ContextCriteria returns the criteria for the supplied context.
func ContextCriteria(context *ctx.Context) *Criteria {
	return context.Get(criteriaContextKey).(*Criteria)
}
