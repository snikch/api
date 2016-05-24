package vc

import "net/http"

// Response represents a non error display response.
type Response struct {
	Payload interface{} `json:"payload"`
	// A pointer is used here to allow empty maps to be returned.
	Sideload *map[string]map[string]interface{} `json:"related,omitempty"`
}

// RespondWithStatusCode returns an empty response.
func RespondWithStatusCode(w http.ResponseWriter, r *http.Request, code int) {
	w.WriteHeader(code)
}

// RespondWithData will use the default renderer to render the response.
func RespondWithData(w http.ResponseWriter, r *http.Request, data interface{}, code int) {
	RespondWithRenderedData(DefaultRenderer, w, r, data, code)
}

// RespondWithRenderedData will render the supplied data as a response.
func RespondWithRenderedData(renderer Renderer, w http.ResponseWriter, r *http.Request, data interface{}, code int) {
	body, err := renderer.Render(data)
	if err != nil {
		RespondWithError(w, r, err)
		return
	}

	// If the code is empty we assume it's an ok response.
	if code == 0 {
		code = http.StatusOK
	}
	w.WriteHeader(code)
	w.Write(body)
}
