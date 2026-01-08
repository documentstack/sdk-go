// Package documentstack provides a Go SDK for the DocumentStack PDF generation API.
package documentstack

// Config holds configuration options for the DocumentStack client.
type Config struct {
	// APIKey is the API key for authentication (Bearer token). Required.
	APIKey string

	// BaseURL is the base URL of the DocumentStack API.
	// Default: "https://api.documentstack.dev"
	BaseURL string

	// Timeout is the request timeout in seconds.
	// Default: 30
	Timeout int

	// Headers are custom headers to include in all requests.
	Headers map[string]string

	// Debug enables debug logging.
	Debug bool
}

// GenerateOptions contains options for PDF generation.
type GenerateOptions struct {
	// Filename is the custom filename for the generated PDF (without .pdf extension).
	Filename string `json:"filename,omitempty"`
}

// GenerateRequest is the request payload for PDF generation.
type GenerateRequest struct {
	// Data is the template data for variable substitution.
	Data map[string]interface{} `json:"data,omitempty"`

	// Options contains generation options.
	Options *GenerateOptions `json:"options,omitempty"`
}

// GenerateResponse contains the generated PDF and metadata.
type GenerateResponse struct {
	// PDF is the PDF binary data.
	PDF []byte

	// Filename is the filename from Content-Disposition header.
	Filename string

	// GenerationTimeMs is the generation time in milliseconds.
	GenerationTimeMs int64

	// ContentLength is the content length in bytes.
	ContentLength int64
}

// APIErrorResponse represents an error response from the API.
type APIErrorResponse struct {
	Error   string      `json:"error"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}
