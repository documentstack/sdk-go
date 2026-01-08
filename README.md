# documentstack-go

Official Go SDK for the DocumentStack PDF generation API.

## Installation

```bash
go get github.com/documentstack/sdk-go
```

## Quick Start

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/documentstack/sdk-go"
)

func main() {
	// Initialize the client
	client, err := documentstack.New(documentstack.Config{
		APIKey: "your-api-key",
	})
	if err != nil {
		log.Fatal(err)
	}

	// Generate a PDF
	result, err := client.Generate(context.Background(), "template-id", &documentstack.GenerateRequest{
		Data: map[string]interface{}{
			"name":   "John Doe",
			"amount": 1500,
		},
		Options: &documentstack.GenerateOptions{
			Filename: "invoice",
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	// Save to file
	if err := os.WriteFile(result.Filename, result.PDF, 0644); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("PDF generated in %dms\n", result.GenerationTimeMs)
}
```

## Configuration

```go
client, err := documentstack.New(documentstack.Config{
	// Required: Your API key
	APIKey: "your-api-key",

	// Optional: Custom base URL (default: https://api.documentstack.dev)
	BaseURL: "https://api.documentstack.dev",

	// Optional: Request timeout in seconds (default: 30)
	Timeout: 30,

	// Optional: Custom headers for all requests
	Headers: map[string]string{
		"X-Custom-Header": "value",
	},

	// Optional: Enable debug logging (default: false)
	Debug: false,
})
```

## API Reference

### `client.Generate(ctx, templateID, request)`

Generate a PDF from a template.

**Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `ctx` | `context.Context` | Yes | Context for cancellation |
| `templateID` | `string` | Yes | The ID of the template to use |
| `request.Data` | `map[string]interface{}` | No | Template data |
| `request.Options.Filename` | `string` | No | Custom filename |

**Returns:** `*GenerateResponse, error`

```go
type GenerateResponse struct {
	PDF              []byte // PDF binary data
	Filename         string // Filename from response
	GenerationTimeMs int64  // Generation time in ms
	ContentLength    int64  // File size in bytes
}
```

## Error Handling

The SDK provides typed errors for different failure scenarios:

```go
import "github.com/documentstack/sdk-go"

result, err := client.Generate(ctx, "template-id", nil)
if err != nil {
	if apiErr, ok := err.(*documentstack.APIError); ok {
		switch {
		case apiErr.IsValidationError():
			// 400: Invalid request
			fmt.Printf("Validation failed: %s\n", apiErr.Message)
		case apiErr.IsAuthenticationError():
			// 401: Invalid API key
			fmt.Printf("Authentication failed: %s\n", apiErr.Message)
		case apiErr.IsForbiddenError():
			// 403: No access to template
			fmt.Printf("Access forbidden: %s\n", apiErr.Message)
		case apiErr.IsNotFoundError():
			// 404: Template not found
			fmt.Printf("Template not found: %s\n", apiErr.Message)
		case apiErr.IsRateLimitError():
			// 429: Rate limit exceeded
			if rlErr, ok := err.(*documentstack.RateLimitError); ok {
				fmt.Printf("Rate limited. Retry after %d seconds\n", rlErr.RetryAfter)
			}
		case apiErr.IsServerError():
			// 500: Server error
			fmt.Printf("Server error: %s\n", apiErr.Message)
		}
	} else if _, ok := err.(*documentstack.TimeoutError); ok {
		fmt.Println("Request timed out")
	} else if netErr, ok := err.(*documentstack.NetworkError); ok {
		fmt.Printf("Network error: %s\n", netErr.Message)
	}
}
```

## Context Support

The SDK fully supports Go contexts for cancellation and timeouts:

```go
// With timeout
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

result, err := client.Generate(ctx, "template-id", request)

// With cancellation
ctx, cancel := context.WithCancel(context.Background())
go func() {
	// Cancel after some condition
	cancel()
}()

result, err := client.Generate(ctx, "template-id", request)
```

## Requirements

- Go 1.21 or higher
- No external dependencies (uses only standard library)

## License

MIT
