package fail

import "net/http"

// BadRequestError represents a bad request error.
type BadRequestError struct {
	Err
}

// NewBadRequestError returns a new BadRequestError to wrap the supplied error.
func NewBadRequestError(err error) BadRequestError {
	return BadRequestError{
		Err: Err{
			OriginalError: err,
		},
	}
}

// StatusCode implements the `vc.StatusError` interface.
func (err BadRequestError) StatusCode() int {
	return http.StatusBadRequest
}
