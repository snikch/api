package middleware

import "net/http"

// Middleware represents a simple chained handler loop.
type Middleware []http.Handler

// New returns a new middleware instance.
func New() Middleware {
	return Middleware{}
}

// Add takes an array of handlers and appends them to the chain.
func (m *Middleware) Add(handler ...http.Handler) {
	*m = append(*m, handler...)
}

// ServeHTTP implements http.Handler and loops over the handlers until one
// writes to the response.
func (m Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mw := NewResponseWriter(w)
	for _, handler := range m {
		if mw.written {
			return
		}
		handler.ServeHTTP(mw, r)
	}
}

// ResponseWriter makes a write aware http.ResponseWriter.
type ResponseWriter struct {
	http.ResponseWriter
	written bool
}

// NewResponseWriter takes a response writer and wraps it.
func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{
		ResponseWriter: w,
	}
}

// Write defers to the underlying response writer but marks the response as
// written.
func (w *ResponseWriter) Write(bytes []byte) (int, error) {
	w.written = true
	return w.ResponseWriter.Write(bytes)
}

// WriteHeader defers to the underlying response writer but marks the response
// as written.
func (w *ResponseWriter) WriteHeader(code int) {
	w.written = true
	w.ResponseWriter.WriteHeader(code)
}
