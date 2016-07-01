package fail

import (
	"database/sql"
	"net/http"
)

// PossibleRowNotFoundError represents a resource not found.
type PossibleRowNotFoundError struct {
	Err
}

// NewPossibleRowNotFoundError determines if err matches sql.ErrNoRows and
// returns the err wrapped in a NotFoundError, otherwise the err returned
// as normal.
func NewPossibleRowNotFoundError(err error) error {
	if err == sql.ErrNoRows {
		return NotFoundError{
			Err: Err{
				OriginalError: err,
				Description:   "The requested resource cannot be found",
			},
		}
	}
	return err
}

// StatusCode implements the `vc.StatusError` interface.
func (err PossibleRowNotFoundError) StatusCode() int {
	return http.StatusNotFound
}
