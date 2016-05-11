package fail

// Err reprents a type that contains fields for additional error metadata.
// This type is generally extended with more domain specific context.
type Err struct {
	Code             int
	OriginalError    error
	Description      string
	AdditionalFields map[string]string
}

// Error implements the `error` interface.
func (err Err) Error() string {
	return err.OriginalError.Error()
}

// ErrorCode implements the `vc.ErrorCode` interface.
func (err Err) ErrorCode() int {
	return err.Code
}

// ErrorDescription implements `vc.DescriptiveError` and should return a human
// readable description of the error.
func (err Err) ErrorDescription() string {
	return err.Description
}

// ErrorFields implements `vc.AnnotatedError` and should return a map of error
// details if appropriate.
func (err Err) ErrorFields() map[string]string {
	return err.AdditionalFields
}

// WithField is a chainable method that adds additional fields to the error.
func (err *Err) WithField(key, value string) {
	if err.AdditionalFields == nil {
		err.AdditionalFields = map[string]string{}
	}
	err.AdditionalFields[key] = value
}

// WithFields is a chainable method that adds additional fields to the error.
func (err *Err) WithFields(fields map[string]string) {
	if err.AdditionalFields == nil {
		err.AdditionalFields = map[string]string{}
	}
	for key, value := range fields {
		err.AdditionalFields[key] = value
	}
}
