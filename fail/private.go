package fail

import (
	"crypto/sha1"
	"fmt"
	"time"
)

// Private masks an error with a public display message and tracking id to hide
// any sensitive system errors, and allow tracking of errors.
// Consumers should use the tracking id to link the public error to the private
// error message which should be logged into internal systems only.
type Private struct {
	ID            string
	PublicMessage string
	OriginalErr   error
}

// NewPrivate returns a private error wrapping the supplied error.
func NewPrivate(err error) *Private {
	// Generate a randomish 16 character id by taking the nano time and sha'ing it.
	// This definitely isn't a perfect solution, but really, the chances of multiple
	// errors at the same nanosecond are pretty slim - so collisions are unlikely.
	timeSHA := sha1.Sum([]byte(fmt.Sprintf("%d", time.Now().UnixNano())))
	id := fmt.Sprintf("%x", timeSHA[:])[:16]
	return &Private{
		ID:            id,
		PublicMessage: fmt.Sprintf("An unexpected error occurred [%s]", id),
		OriginalErr:   err,
	}
}

// Error implements the error interface.
func (private Private) Error() string {
	return private.PublicMessage
}

// ErrorFields implements vc.AnnotatedError and allows for the structured
// logging of the tracking id, and original private errors.
func (private Private) ErrorFields() map[string]string {
	return map[string]string{
		"id":             private.ID,
		"original_error": private.OriginalErr.Error(),
	}
}
