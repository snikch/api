package fail

// ValidationErrorStatusCode represents the HTTP status code for ValidationError.
var ValidationErrorStatusCode = 422

// ValidationError represents a validation error.
type ValidationError struct {
	Err
}

// NewValidationError returns a new ValidationError to wrap the supplied error.
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
