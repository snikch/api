package fail

import (
	"fmt"
	"net/http"
)

// AuthenticationError represents a failure to authenticate.
type AuthenticationError struct {
	Err
}

// NewAuthError returns a new AuthenticationError to wrap the supplied error.
func NewAuthError(code int, parts ...string) AuthenticationError {
	var description, err string

	if len(parts) >= 1 {
		err = parts[0]
	}
	if len(parts) >= 2 {
		description = parts[1]
	}

	return AuthenticationError{
		Err: Err{
			Code:          code,
			OriginalError: fmt.Errorf(err),
			Description:   description,
		},
	}
}

// StatusCode implements the `vc.StatusError` interface. Any
// AuthenticationError should return a 401 response.
func (err AuthenticationError) StatusCode() int {
	return http.StatusUnauthorized
}
