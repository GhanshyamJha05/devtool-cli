package utils

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetchData(t *testing.T) {
	// Create a mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "success"}`))
	}))
	defer server.Close()

	// Test valid URL
	resp, err := FetchData(server.URL)
	if err != nil {
		t.Fatalf("FetchData failed for valid URL: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", resp.StatusCode)
	}
	if resp.Body != `{"message": "success"}` {
		t.Errorf("Unexpected response body: %s", resp.Body)
	}

	// Test invalid URL format
	_, err = FetchData("not-a-url")
	if err == nil {
		t.Errorf("Expected error for invalid URL, got nil")
	}

	// Test unsupported scheme
	_, err = FetchData("ftp://example.com")
	if err == nil {
		t.Errorf("Expected error for unsupported scheme, got nil")
	}
}
