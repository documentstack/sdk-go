// Package documentstack provides a Go SDK for the DocumentStack PDF generation API.
package documentstack

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	defaultBaseURL = "https://api.documentstack.dev"
	defaultTimeout = 30
)

// Client is the DocumentStack API client.
type Client struct {
	config     Config
	httpClient *http.Client
}

// New creates a new DocumentStack client with the given configuration.
//
// Example:
//
//	client, err := documentstack.New(documentstack.Config{
//		APIKey: "your-api-key",
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	result, err := client.Generate(context.Background(), "template-id", &documentstack.GenerateRequest{
//		Data: map[string]interface{}{
//			"name":   "John Doe",
//			"amount": 100,
//		},
//		Options: &documentstack.GenerateOptions{
//			Filename: "invoice",
//		},
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	os.WriteFile("invoice.pdf", result.PDF, 0644)
func New(config Config) (*Client, error) {
	if config.APIKey == "" {
		return nil, &DocumentStackError{Message: "API key is required"}
	}

	if config.BaseURL == "" {
		config.BaseURL = defaultBaseURL
	}
	config.BaseURL = strings.TrimSuffix(config.BaseURL, "/")

	if config.Timeout <= 0 {
		config.Timeout = defaultTimeout
	}

	if config.Headers == nil {
		config.Headers = make(map[string]string)
	}

	client := &Client{
		config: config,
		httpClient: &http.Client{
			Timeout: time.Duration(config.Timeout) * time.Second,
		},
	}

	return client, nil
}

// Generate generates a PDF from a template.
//
// Parameters:
//   - ctx: Context for cancellation and timeouts
//   - templateID: The ID of the template to use
//   - request: Generation request with data and options (can be nil)
//
// Returns the generated PDF and metadata, or an error.
func (c *Client) Generate(ctx context.Context, templateID string, request *GenerateRequest) (*GenerateResponse, error) {
	if templateID == "" {
		return nil, NewValidationError("Template ID is required", nil)
	}

	if request == nil {
		request = &GenerateRequest{}
	}

	endpoint := fmt.Sprintf("%s/api/v1/generate/%s", c.config.BaseURL, url.PathEscape(templateID))

	body, err := json.Marshal(request)
	if err != nil {
		return nil, &NetworkError{Message: "failed to marshal request body", Cause: err}
	}

	if c.config.Debug {
		log.Printf("[DocumentStack] Request: POST %s\n", endpoint)
		log.Printf("[DocumentStack] Body: %s\n", string(body))
	}

	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, &NetworkError{Message: "failed to create request", Cause: err}
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.config.APIKey))

	for key, value := range c.config.Headers {
		req.Header.Set(key, value)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return nil, &TimeoutError{Timeout: c.config.Timeout}
		}
		return nil, &NetworkError{Message: "request failed", Cause: err}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, c.parseErrorResponse(resp)
	}

	pdf, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &NetworkError{Message: "failed to read response body", Cause: err}
	}

	// Extract metadata from headers
	contentDisposition := resp.Header.Get("Content-Disposition")
	generationTimeMs, _ := strconv.ParseInt(resp.Header.Get("X-Generation-Time-Ms"), 10, 64)
	contentLength, _ := strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 64)

	// Parse filename from Content-Disposition
	filename := "document.pdf"
	re := regexp.MustCompile(`filename="?([^";\n]+)"?`)
	if matches := re.FindStringSubmatch(contentDisposition); len(matches) > 1 {
		filename = matches[1]
	}

	if contentLength == 0 {
		contentLength = int64(len(pdf))
	}

	if c.config.Debug {
		log.Printf("[DocumentStack] Response: filename=%s, time=%dms, size=%d\n", filename, generationTimeMs, contentLength)
	}

	return &GenerateResponse{
		PDF:              pdf,
		Filename:         filename,
		GenerationTimeMs: generationTimeMs,
		ContentLength:    contentLength,
	}, nil
}

// parseErrorResponse parses an error response from the API.
func (c *Client) parseErrorResponse(resp *http.Response) error {
	var errorBody APIErrorResponse

	body, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &errorBody); err != nil {
		errorBody = APIErrorResponse{
			Error:   "Unknown Error",
			Message: resp.Status,
		}
	}

	apiErr := &APIError{
		StatusCode: resp.StatusCode,
		ErrorCode:  errorBody.Error,
		Message:    errorBody.Message,
		Details:    errorBody.Details,
	}

	if resp.StatusCode == 429 {
		retryAfter, _ := strconv.Atoi(resp.Header.Get("Retry-After"))
		return &RateLimitError{
			APIError:   apiErr,
			RetryAfter: retryAfter,
		}
	}

	return apiErr
}
