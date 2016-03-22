package fail

import "net/http"

type BadRequestError struct {
	Err
}

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
