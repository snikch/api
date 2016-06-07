package metrics

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/rcrowley/go-metrics"
	"github.com/snikch/api/middleware"
)

// MeterByStatus is an http.Handler that counts responses by their HTTP
// status code via go-metrics.
// This code is ripped blatantly from tigertonic, and converted to meters.
type MeterByStatus struct {
	sync.Mutex
	meters   map[int]metrics.Meter
	handler  http.Handler
	name     string
	registry metrics.Registry
}

// MeteredByStatus returns an http.Handler that passes requests to an
// underlying http.Handler and then counts the response by its HTTP status code
// via go-metrics.
func MeteredByStatus(
	handler http.Handler,
	name string,
	registry metrics.Registry,
) *MeterByStatus {
	if nil == registry {
		registry = metrics.DefaultRegistry
	}
	return &MeterByStatus{
		meters:   map[int]metrics.Meter{},
		handler:  handler,
		name:     name,
		registry: registry,
	}
}

// ServeHTTP passes the request to the underlying http.Handler and then counts
// the response by its HTTP status code via go-metrics.
func (c *MeterByStatus) ServeHTTP(w0 http.ResponseWriter, r *http.Request) {
	w := middleware.NewResponseWriter(w0)
	c.handler.ServeHTTP(w, r)
	c.Lock()
	_, ok := c.meters[w.StatusCode]
	if !ok {
		meter := metrics.NewMeter()
		c.meters[w.StatusCode] = meter
		if err := c.registry.Register(
			fmt.Sprintf("meter-%s-%d", c.name, w.StatusCode),
			meter,
		); nil != err {
			c.Unlock()
			panic(err)
		}
	}
	c.Unlock()
	c.meters[w.StatusCode].Mark(1)
}

// TimedByStatus is an http.Handler that times responses by their HTTP status
// code via go-metrics.
// This code is ripped blatantly from tigertonic, and converted to meters.
type TimedByStatus struct {
	sync.Mutex
	timers   map[int]metrics.Timer
	handler  http.Handler
	name     string
	registry metrics.Registry
}

// NewTimedByStatus returns an http.Handler that passes requests to an
// underlying http.Handler and then times the response by its HTTP status code
// via go-metrics.
func NewTimedByStatus(
	handler http.Handler,
	name string,
	registry metrics.Registry,
) *TimedByStatus {
	if nil == registry {
		registry = metrics.DefaultRegistry
	}
	return &TimedByStatus{
		timers:   map[int]metrics.Timer{},
		handler:  handler,
		name:     name,
		registry: registry,
	}
}

// ServeHTTP passes the request to the underlying http.Handler and then times
// the response by its HTTP status code via go-metrics.
func (c *TimedByStatus) ServeHTTP(w0 http.ResponseWriter, r *http.Request) {
	w := middleware.NewResponseWriter(w0)
	// Mark the start time.
	start := time.Now()
	// Run the request.
	c.handler.ServeHTTP(w, r)
	c.Lock()
	_, ok := c.timers[w.StatusCode]
	if !ok {
		// Generate a new timer if required.
		timer := metrics.NewTimer()
		c.timers[w.StatusCode] = timer
		if err := c.registry.Register(
			fmt.Sprintf("timer-%s-%d", c.name, w.StatusCode),
			timer,
		); nil != err {
			c.Unlock()
			panic(err)
		}
	}
	c.Unlock()
	c.timers[w.StatusCode].UpdateSince(start)
}
