package metrics

import (
	"fmt"
	"net/http"

	"github.com/rcrowley/go-metrics"
	"github.com/snikch/api/middleware"
)

// MeterByStatus is an http.Handler that counts responses by their HTTP
// status code via go-metrics.
// This code is ripped blatantly from tigertonic, and converted to meters.
type MeterByStatus struct {
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
	_, ok := c.meters[w.StatusCode]
	if !ok {
		meter := metrics.NewMeter()
		c.meters[w.StatusCode] = meter
		if err := c.registry.Register(
			fmt.Sprintf("%s-%d", c.name, w.StatusCode),
			meter,
		); nil != err {
			panic(err)
		}
	}
	c.meters[w.StatusCode].Mark(1)
}
