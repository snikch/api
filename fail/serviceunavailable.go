package fail

import "net/http"

// ServiceUnavailable represents a service unavailable error.
type ServiceUnavailable struct {
	Err
}

// NewServiceUnavailableError returns a new ServiceUnavailable to wrap the supplied error.
func NewServiceUnavailableError(err error) ServiceUnavailable {
	return ServiceUnavailable{
		Err: Err{
			OriginalError: err,
		},
	}
}

// StatusCode implements the `vc.StatusError` interface.
func (err ServiceUnavailable) StatusCode() int {
	return http.StatusServiceUnavailable
}
