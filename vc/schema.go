package vc

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/snikch/api/ctx"
	"github.com/snikch/api/fail"
	schema "github.com/xeipuuv/gojsonschema"
)

func MustSchema(schemaString string) *schema.Schema {
	s, err := schema.NewSchema(schema.NewStringLoader(schemaString))
	if err != nil {
		panic(err)
	}
	return s
}

func MustSchemaFromFile(file string) *schema.Schema {
	contents, err := ioutil.ReadFile(file)
	if err != nil {
		panic(fmt.Errorf("schema: %s", err.Error()))
	}

	return MustSchema(string(contents))
}

// UnmarshalAndValidateSchema attempts to validate the body on the suppled
// context against the supplied schema, then unmarshal it into the supplied obj.
func UnmarshalAndValidateSchema(context *ctx.Context, s *schema.Schema, obj interface{}) error {
	// Check we can get a body.
	if context.Request == nil || context.Request.Body == nil {
		return errors.New("No request, or request body, available on context to unmarshal")
	}

	// Read the full body.
	body, err := ioutil.ReadAll(context.Request.Body)
	if err != nil {
		return err
	}

	// Validate the body against the schema.
	result, err := s.Validate(schema.NewStringLoader(string(body)))
	if err != nil {
		return fail.NewBadRequestError(err)
	}

	// Return early if the json is invalid.
	if !result.Valid() {
		err := fail.NewSchemaValidationError(result.Errors())
		err.Description = "It looks like one or more fields weren’t present, of the correct type, or contained an incorrect value. If it’s not obvious what the problem was, make sure you consult the API documentation. If that still doesn’t make it obvious, our documentation clearly isn’t up to scratch. Contact support and we’ll help you out."
		return err
	}

	// Unmarshal the body into the supplied object.
	err = json.Unmarshal(body, obj)
	if err != nil {
		return fail.NewBadRequestError(err)
	}
	return nil
}
