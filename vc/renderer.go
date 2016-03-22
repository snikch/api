package vc

var (
	// DefaultRenderer is the renderer used by default. This can be changed at
	// runtime to a renderer that suits.
	// DefaultRenderer = render.JSONRenderer{}
	DefaultRenderer Renderer
)

type Renderer interface {
	Render(interface{}) ([]byte, error)
	RenderError(APIError) []byte
}
