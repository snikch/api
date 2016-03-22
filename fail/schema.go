package fail

import (
	"errors"
	"fmt"
	"strings"

	schema "github.com/xeipuuv/gojsonschema"
)

// NewSchemaValidationError returns a ValidationError that contains information
// about every field that was invalid.
func NewSchemaValidationError(schemaErrors []schema.ResultError) ValidationError {
	fields := map[string]string{}
	for _, err := range schemaErrors {
		field := err.Field()
		prefix := err.Context().String()

		// If this isn't a root node, i.e. is a sub entity, update the field name.
		if prefix != "(root)" && err.Type() != "invalid_type" {
			field = fmt.Sprintf("%s.%s",
				strings.Replace(prefix, "(root).", "", -1),
				field,
			)
		}

		fields[field] = err.Description()
	}
	err := NewValidationError(errors.New("Invalid data supplied"))
	err.AdditionalFields = fields
	return err
}
