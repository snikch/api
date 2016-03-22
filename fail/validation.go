package fail

var ValidationErrorStatusCode = 422

type ValidationError struct {
	Err
}

func NewValidationError(err error) ValidationError {
	return ValidationError{
		Err: Err{
			OriginalError: err,
		},
	}
}

// StatusCode implements the `vc.StatusError` interface. Any
// AuthenticationError should return a 401 response.
func (err ValidationError) StatusCode() int {
	return ValidationErrorStatusCode
}
