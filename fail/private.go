package fail

import (
	"crypto/sha1"
	"fmt"
	"net/http"
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
	Code          int
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
		PublicMessage: "An unexpected error occurred",
		OriginalErr:   err,
		Code:          http.StatusInternalServerError,
	}
}

// WithMessage is a chainable method to set the public message.
func (private *Private) WithMessage(message string) *Private {
	private.PublicMessage = message
	return private
}

// WithStatusCode is a chainable method to set the status code.
func (private *Private) WithStatusCode(code int) *Private {
	private.Code = code
	return private
}

// StatusCode implements the StatusError interface.
func (private Private) StatusCode() int {
	return private.Code
}

// Error implements the error interface.
func (private Private) Error() string {
	return fmt.Sprintf("[%s] %s", private.ID, private.PublicMessage)
}

// LogFields implements vc.StructuredLogsError and allows for the structured
// logging of the tracking id, and original private errors.
func (private Private) LogFields() map[string]string {
	return map[string]string{
		"id":             private.ID,
		"original_error": private.OriginalErr.Error(),
	}
}
