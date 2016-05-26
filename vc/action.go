package vc

import (
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/rcrowley/go-metrics"
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
	MetricsRegistry metrics.Registry
}

func NewActionProcessor() *ActionProcessor {
	return &ActionProcessor{
		MetricsRegistry: metrics.NewRegistry(),
	}
}

// ActionHandler implementers are responsible for returning payload data for
// a request, along with a status code or error.
type ActionHandler interface {
	HandleAction(*ctx.Context) (interface{}, int, error)
}

// ActionHandlerFunc wraps a function with the HandleAction signature to a full
// ActionHandler interface.
type ActionHandlerFunc struct {
	Handler func(*ctx.Context) (interface{}, int, error)
}

// HandleAction implements the ActionHander interface and simply calls the
// underlying function.
func (fn ActionHandlerFunc) HandleAction(context *ctx.Context) (interface{}, int, error) {
	return fn.Handler(context)
}

// HandleActionFunc returns an http.Handler for the suppled action function.
// A type and action name are used in metrics and logging functions.
func (p *ActionProcessor) HandleActionFunc(typ, action string, fn func(*ctx.Context) (interface{}, int, error)) httprouter.Handle {
	return p.HTTPHandler(typ, action, ActionHandlerFunc{
		Handler: fn,
	})
}

var requestCriteriaTransformers = []func(*ctx.Context, *Criteria){}

func RegisterCriteriaTransformer(transformer func(*ctx.Context, *Criteria)) {
	requestCriteriaTransformers = append(requestCriteriaTransformers, transformer)
}

// HTTPHandler takes an ActionHandler and returns a http.Handler instance
// that can be used. The type and action are used to determine the context in
// several areas, such as transformers and metrics.
func (p *ActionProcessor) HTTPHandler(typ, action string, handler ActionHandler) httprouter.Handle {
	// Create a new timer for timing this handler.
	timer := metrics.NewTimer()
	p.MetricsRegistry.Register(typ+"-"+action, timer)
	sideloadTimer := metrics.NewTimer()
	p.MetricsRegistry.Register(typ+"-"+action+"-sideload", sideloadTimer)
	unlockTimer := metrics.NewTimer()
	p.MetricsRegistry.Register(typ+"-"+action+"-unlock", unlockTimer)

	return httprouter.Handle(func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		// At the end of this function, add a time metric.
		defer timer.UpdateSince(time.Now())

		// Create a new context for this action.
		context := ctx.NewContext()
		context.Request = r
		context.EntityType = typ
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
			start := time.Now()
			// Retrieve any sideloaded entities.
			sideloaded, err := sideload.Load(context, payload, criteria.Sideload)

			response.Sideload = &sideloaded
			if err != nil {
				timer.UpdateSince(start)
				RespondWithError(w, r, err)
				return
			}
			timer.UpdateSince(start)
		}

		// Unlock any entities registered for this request.
		unlockStartTime := time.Now()
		err = lynx.ContextStore(context).Unlock()
		unlockTimer.UpdateSince(unlockStartTime)

		if err != nil {
			RespondWithError(w, r, err)
			return
		}

		RespondWithData(w, r, response, code)
	})
}
