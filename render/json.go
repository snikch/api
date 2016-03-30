package render

import (
	"encoding/json"

	"github.com/snikch/api/vc"
)

// JSONRenderer is a type that implements the vc.Renderer interface.
type JSONRenderer struct {
}

// Render marshals the supplied data to json.
func (j JSONRenderer) Render(data interface{}) ([]byte, error) {
	// A possibly unpopular decision, but if the data is nil, return nothing.
	if data == nil {
		return nil, nil
	}
	return json.MarshalIndent(data, "", "  ")
}

type jsonError struct {
	Error       string            `json:"error"`
	Description string            `json:"description,omitempty"`
	Code        int               `json:"code,omitempty"`
	Fields      map[string]string `json:"fields,omitempty"`
}

func (j JSONRenderer) RenderError(error vc.APIError) []byte {
	marshalled, err := json.MarshalIndent(jsonError{
		Error:       error.Error,
		Description: error.Description,
		Code:        error.Code,
		Fields:      error.Fields,
	}, "", "  ")
	if err != nil {
		return []byte(err.Error())
	}
	return marshalled
}
