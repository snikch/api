package vc

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/snikch/api/ctx"
	"github.com/snikch/api/lynx"
	"github.com/snikch/api/sideload"
)

// EmptyResponse is used to determine if a response should be empty.
var EmptyResponse = &Response{}

type contextKey int

const (
	criteriaContextKey contextKey = iota
	paramsContextKey
)

// ActionProcessor handles an entire action lifecycle, from data retrieval
// through to unlocking, and responding.
type ActionProcessor struct {
	SideloadEnabled bool
}

// ActionHandler implementers are responsible for returning payload data for
// a request, along with a status code or error.
type ActionHandler interface {
	// ActionEntityType can be used to return a classification for the request.
	// This can then be used in transformers and other context agnostic areas.
	ActionEntityType() string
	HandleAction(*ctx.Context) (interface{}, int, error)
}

// ActionHandlerFunc wraps a function with the HandleAction signature to a full
// ActionHandler interface.
type ActionHandlerFunc struct {
	Handler func(*ctx.Context) (interface{}, int, error)
	Type    string
}

// HandleAction implements the ActionHander interface and simply calls the
// underlying function.
func (fn ActionHandlerFunc) HandleAction(context *ctx.Context) (interface{}, int, error) {
	return fn.Handler(context)
}

func (fn ActionHandlerFunc) ActionEntityType() string {
	return fn.Type
}

// HandleActionFunc returns an http.Handler for the suppled action function.
func (p *ActionProcessor) HandleActionFunc(typ string, fn func(*ctx.Context) (interface{}, int, error)) httprouter.Handle {
	return p.HTTPHandler(ActionHandlerFunc{
		Handler: fn,
		Type:    typ,
	})
}

var requestCriteriaTransformers = []func(*ctx.Context, *Criteria){}

func RegisterCriteriaTransformer(transformer func(*ctx.Context, *Criteria)) {
	requestCriteriaTransformers = append(requestCriteriaTransformers, transformer)
}

// HTTPHandler takes an ActionHandler and returns a http.Handler instance
// that can be used.
func (p *ActionProcessor) HTTPHandler(handler ActionHandler) httprouter.Handle {
	return httprouter.Handle(func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		// Create a new context for this action.
		context := ctx.NewContext()
		context.Request = r
		context.EntityType = handler.ActionEntityType()
		SetContextParams(context, params)

		// Get any criteria, and transform it if required.
		criteria := RequestCriteria(r)
		for _, transformer := range requestCriteriaTransformers {
			transformer(context, criteria)
		}
		// Make the criteria available on the content.
		SetContextCriteria(context, criteria)

		// Get the base payload back from the ActionHandler instance.
		payload, code, err := handler.HandleAction(context)
		if err != nil {
			RespondWithError(w, r, err)
			return
		}

		// If an empty response is required, return an empty response.
		if payload == EmptyResponse {
			RespondWithStatusCode(w, r, code)
			return
		}

		// Build up a response.
		response := Response{
			Payload: payload,
		}

		if p.SideloadEnabled {
			// Retrieve any sideloaded entities.
			sideloaded, err := sideload.Load(context, payload, criteria.Sideload)

			response.Sideload = &sideloaded
			if err != nil {
				RespondWithError(w, r, err)
				return
			}
		}

		// Unlock any entities registered for this request.
		err = lynx.ContextStore(context).Unlock()
		if err != nil {
			RespondWithError(w, r, err)
			return
		}

		RespondWithData(w, r, response, code)
	})
}
