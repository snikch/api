package vc

import (
	"github.com/julienschmidt/httprouter"
	"github.com/snikch/api/ctx"
)

// SetContextParams sets the supplied httprouter.Params on the context.
func SetContextParams(context *ctx.Context, params httprouter.Params) {
	context.Set(paramsContextKey, params)
}

// ContextParams returns the httprouter.Params for the context.
func ContextParams(context *ctx.Context) httprouter.Params {
	return context.Get(paramsContextKey).(httprouter.Params)
}
