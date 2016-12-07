package fail

import (
	"fmt"
	"net/http"
	"strconv"
)

// PermissionsError represents a forbidden access request.
type PermissionsError struct {
	Err
}

// NewPermissionsError returns a new PermissionsError to wrap the supplied error.
// The parts slice elements refer to the original error and a more detailed
// description respectively.
func NewPermissionsError(code int, parts ...string) PermissionsError {
	var description, err string

	if len(parts) >= 1 {
		err = parts[0]
	}
	if len(parts) >= 2 {
		description = parts[1]
	}

	// Map our custom permission error code.
	internalErrCode := map[string]string{
		"error_code": strconv.Itoa(code),
	}

	return PermissionsError{
		Err: Err{
			Code:             code,
			OriginalError:    fmt.Errorf(err),
			Description:      description,
			AdditionalFields: internalErrCode,
		},
	}
}

// StatusCode implements the `vc.StatusError` interface.
func (err PermissionsError) StatusCode() int {
	return http.StatusForbidden
}
