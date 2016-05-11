package vc

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/snikch/api/fail"
	"github.com/snikch/api/log"
)

// StatusError defines an interface for an error that also includes a custom
// http status code. Use this liberally to differentiate service errors from
// things like authentication, authorization or validation errors.
type StatusError interface {
	StatusCode() int
}

// DescriptiveError defines an interface for returning a human readable error
// description. It’s ok to use humour here, as these descriptions will be read
// by developers, or at least real people, and can often be read in a time of
// frustration. Be authentic with your descriptions, and good karma may come
// your way.
type DescriptiveError interface {
	ErrorDescription() string
}

// AnnotatedError defines an interface for additional metadata for an error.
// This is where you want to put all of your error data when it's more than just
// a single error case condition, for example all of the invalid fields in a bad
// update or create operation.
type AnnotatedError interface {
	ErrorFields() map[string]string
}

// StructuredLogsError defines an interface for additional metadata for an error.
// This metadata will be logged, but not output to the public error response.
type StructuredLogsError interface {
	LogFields() map[string]string
}

// APIError represents an API error response. This will be passed to a renderer
// error method for conversion into an appropriate response.
type APIError struct {
	Code        int               `json:"code,omitempty"`
	Description string            `json:"description,omitempty"`
	Error       string            `json:"error"`
	Fields      map[string]string `json:"fields,omitempty"`
}

// RespondWithError will return an error response with the appropriate message,
// and status codes set.
func RespondWithError(w http.ResponseWriter, r *http.Request, err error) {
	isPublicError := false
	errorResponse := APIError{
		Error: err.Error(),
	}
	code := http.StatusInternalServerError
	if statusErr, ok := err.(StatusError); ok {
		// If we get a status code, this error can be considered a public error.
		isPublicError = true
		code = statusErr.StatusCode()
	}

	if descriptiveErr, ok := err.(DescriptiveError); ok {
		errorResponse.Description = descriptiveErr.ErrorDescription()
	}

	if annotatedErr, ok := err.(AnnotatedError); ok {
		errorResponse.Fields = annotatedErr.ErrorFields()
	}

	w.WriteHeader(code)

	// Now log some information about the failure.
	logData := map[string]interface{}{}
	if structuredLogErr, ok := err.(StructuredLogsError); ok {
		for key, value := range structuredLogErr.LogFields() {
			logData[key] = value
		}
	}
	logData["status"] = code

	if !isPublicError {
		logData["original_error"] = err.Error()
		err = fail.NewPrivate(err)
		errorResponse.Error = err.Error()
	}

	body := DefaultRenderer.RenderError(errorResponse)
	w.Write(body)

	log.WithError(err).WithFields(logrus.Fields(logData)).Error("Returning error response")

}
