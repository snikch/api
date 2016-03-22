package log

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/snikch/api/vc"
)

// AccessLogger wraps an http.Handler and provides request level logging.
type AccessLogger struct {
	http.Handler
}

func NewAccessLogger(handler http.Handler) *AccessLogger {
	return &AccessLogger{
		Handler: handler,
	}
}

// ServeHTTP implements the http.Handler interface and will record information
// about a request, and log it after the request runs.
func (logger AccessLogger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Retrieve the last client ip from the RemoteAddr header field.
	clientIP := r.RemoteAddr
	if colon := strings.LastIndex(clientIP, ":"); colon != -1 {
		clientIP = clientIP[:colon]
	}

	// Create an access log record.
	record := &AccessLogRecord{
		ResponseWriter: w,
		ip:             clientIP,
		time:           time.Now(),
		method:         r.Method,
		uri:            r.RequestURI,
		protocol:       r.Proto,
	}

	// Set up a function to run once the request has been served.
	defer func(record *AccessLogRecord, startTime time.Time) {
		// Recover from a panic if possible.
		if recovered := recover(); recovered != nil {
			var err error
			// Ensure we have an error interface.
			if thisErr, ok := recovered.(error); ok {
				err = thisErr
			} else {
				err = fmt.Errorf("%s", recovered)
			}
			if err != nil {
				WithError(err).Error("Recovered from panic")
				vc.RespondWithError(w, r, err)
			}
		}
		// Log the response info.
		finishTime := time.Now()
		record.time = finishTime.UTC()
		record.duration = finishTime.Sub(startTime)

		WithFields(record.Data()).Infof("")
	}(record, time.Now())

	// Serve the request.
	logger.Handler.ServeHTTP(record, r)
}

// AccessLogRecord is an http.ResponseWriter implementer and wrapper that also
// records information about the request.
type AccessLogRecord struct {
	http.ResponseWriter

	ip                    string
	time                  time.Time
	method, uri, protocol string
	status                int
	length                int64
	duration              time.Duration
}

// Data returns structure log data about a request.
func (r *AccessLogRecord) Data() logrus.Fields {
	return logrus.Fields{
		"ip":       r.ip,
		"finish":   r.time.UnixNano(),
		"method":   r.method,
		"uri":      r.uri,
		"protocol": r.protocol,
		"status":   r.status,
		"length":   r.length,
		"duration": r.duration.String(),
	}
}

// Write implements the http.Handler#Write method, and records the length of
// any data written to the response.
func (r *AccessLogRecord) Write(p []byte) (int, error) {
	written, err := r.ResponseWriter.Write(p)
	r.length += int64(written)
	return written, err
}

// WriteHeader implements the http.Handler#WriteHeader method, and records the
// response status code.
func (r *AccessLogRecord) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}
