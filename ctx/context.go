package ctx

import (
	"net/http"
	"sync"
)

// Context represents a single context whether it be a request or otherwise.
// Passing around a context allows storage of data against the context.
type Context struct {
	Request    *http.Request
	EntityType string
	data       map[interface{}]interface{}
	sync.RWMutex
}

func NewContext() *Context {
	return &Context{
		data: map[interface{}]interface{}{},
	}
}

func (context *Context) Set(key, value interface{}) {
	context.Lock()
	context.SetUnsafe(key, value)
	context.Unlock()
}

func (context *Context) SetUnsafe(key, value interface{}) {
	context.data[key] = value
}

func (context *Context) Get(key interface{}) interface{} {
	value, _ := context.GetOk(key)
	return value
}

func (context *Context) GetOk(key interface{}) (interface{}, bool) {
	context.Lock()
	value, ok := context.GetOkUnsafe(key)
	context.Unlock()
	return value, ok
}

func (context *Context) GetOkUnsafe(key interface{}) (interface{}, bool) {
	value, ok := context.data[key]
	return value, ok
}
