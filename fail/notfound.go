package fail

import "net/http"

type NotFoundError struct {
	Err
}

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
