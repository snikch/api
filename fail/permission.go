package fail

import "net/http"

// PermissionsError represents a forbidden access request.
type PermissionsError struct {
	Err
}

// NewPermissionsError returns a new PermissionsError to wrap the supplied error.
func NewPermissionsError(err error) PermissionsError {
	return PermissionsError{
		Err: Err{
			OriginalError: err,
		},
	}
}

// StatusCode implements the `vc.StatusError` interface.
func (err PermissionsError) StatusCode() int {
	return http.StatusForbidden
}
