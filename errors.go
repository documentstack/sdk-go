package documentstack

import (
	"fmt"
)

// DocumentStackError is the base error type for all SDK errors.
type DocumentStackError struct {
	Message string
}

func (e *DocumentStackError) Error() string {
	return e.Message
}

// APIError is returned when the API request fails with a specific HTTP status.
type APIError struct {
	StatusCode int
	ErrorCode  string
	Message    string
	Details    interface{}
}

func (e *APIError) Error() string {
	return fmt.Sprintf("%s: %s", e.ErrorCode, e.Message)
}

// IsValidationError returns true if the error is a validation error (400).
func (e *APIError) IsValidationError() bool {
	return e.StatusCode == 400
}

// IsAuthenticationError returns true if the error is an authentication error (401).
func (e *APIError) IsAuthenticationError() bool {
	return e.StatusCode == 401
}

// IsForbiddenError returns true if the error is a forbidden error (403).
func (e *APIError) IsForbiddenError() bool {
	return e.StatusCode == 403
}

// IsNotFoundError returns true if the error is a not found error (404).
func (e *APIError) IsNotFoundError() bool {
	return e.StatusCode == 404
}

// IsRateLimitError returns true if the error is a rate limit error (429).
func (e *APIError) IsRateLimitError() bool {
	return e.StatusCode == 429
}

// IsServerError returns true if the error is a server error (5xx).
func (e *APIError) IsServerError() bool {
	return e.StatusCode >= 500
}

// RateLimitError extends APIError with retry information.
type RateLimitError struct {
	*APIError
	RetryAfter int // Seconds to wait before retrying
}

// TimeoutError is returned when a request times out.
type TimeoutError struct {
	Timeout int // Timeout in seconds
}

func (e *TimeoutError) Error() string {
	return fmt.Sprintf("request timed out after %d seconds", e.Timeout)
}

// NetworkError is returned when a network request fails.
type NetworkError struct {
	Message string
	Cause   error
}

func (e *NetworkError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

func (e *NetworkError) Unwrap() error {
	return e.Cause
}

// NewValidationError creates a new validation error.
func NewValidationError(message string, details interface{}) *APIError {
	return &APIError{
		StatusCode: 400,
		ErrorCode:  "Bad Request",
		Message:    message,
		Details:    details,
	}
}

// NewAuthenticationError creates a new authentication error.
func NewAuthenticationError(message string) *APIError {
	return &APIError{
		StatusCode: 401,
		ErrorCode:  "Unauthorized",
		Message:    message,
	}
}

// NewForbiddenError creates a new forbidden error.
func NewForbiddenError(message string) *APIError {
	return &APIError{
		StatusCode: 403,
		ErrorCode:  "Forbidden",
		Message:    message,
	}
}

// NewNotFoundError creates a new not found error.
func NewNotFoundError(message string) *APIError {
	return &APIError{
		StatusCode: 404,
		ErrorCode:  "Not Found",
		Message:    message,
	}
}
