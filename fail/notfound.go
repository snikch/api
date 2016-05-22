package fail

import "net/http"

// NotFoundError represents a resource not found.
type NotFoundError struct {
	Err
}

// NewNotFoundError returns a new NotFoundError to wrap the supplied error.
func NewNotFoundError(err error) NotFoundError {
	return NotFoundError{
		Err: Err{
			OriginalError: err,
		},
	}
}

// StatusCode implements the `vc.StatusError` interface.
func (err NotFoundError) StatusCode() int {
	return http.StatusNotFound
}
