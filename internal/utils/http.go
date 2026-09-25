package utils

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// HTTPResponse holds structured data about the API response.
type HTTPResponse struct {
	StatusCode int
	Status     string
	BodyReader io.ReadCloser
	Duration   time.Duration
	Headers    http.Header
}

// FetchData performs an HTTP GET with timeout, URL validation, and structured response.
func FetchData(rawURL string) (*HTTPResponse, error) {
	// Step 1: Validate the URL before making any network call
	if err := validateURL(rawURL); err != nil {
		return nil, err
	}

	// Step 2: Create a client with a 10-second timeout to avoid hanging forever
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Step 3: Track how long the request takes
	start := time.Now()

	resp, err := client.Get(rawURL)
	if err != nil {
		// Differentiate between common failure types for better UX
		if strings.Contains(err.Error(), "no such host") {
			return nil, fmt.Errorf("DNS resolution failed — host not found: %s", rawURL)
		}
		if strings.Contains(err.Error(), "timeout") {
			return nil, fmt.Errorf("request timed out after 10s: %s", rawURL)
		}
		return nil, fmt.Errorf("network error: %w", err)
	}
	duration := time.Since(start)

	return &HTTPResponse{
		StatusCode: resp.StatusCode,
		Status:     resp.Status,
		BodyReader: resp.Body,
		Duration:   duration,
		Headers:    resp.Header,
	}, nil
}

// validateURL checks that the URL has a valid scheme and host.
func validateURL(rawURL string) error {
	parsed, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return fmt.Errorf("invalid URL format: %s", rawURL)
	}

	if parsed.Scheme == "" || parsed.Host == "" {
		return fmt.Errorf("URL must include scheme and host (e.g., https://example.com), got: %s", rawURL)
	}

	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return fmt.Errorf("unsupported URL scheme '%s' — only http and https are allowed", parsed.Scheme)
	}

	return nil
}
