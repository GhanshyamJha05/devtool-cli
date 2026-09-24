package utils

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFormatJSON(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "devtool-json-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Valid JSON
	validJSONFile := filepath.Join(tempDir, "valid.json")
	os.WriteFile(validJSONFile, []byte(`{"hello":"world"}`), 0644)

	// Invalid JSON
	invalidJSONFile := filepath.Join(tempDir, "invalid.json")
	os.WriteFile(invalidJSONFile, []byte(`{"hello":world"}`), 0644)

	// Empty JSON
	emptyJSONFile := filepath.Join(tempDir, "empty.json")
	os.WriteFile(emptyJSONFile, []byte(` `), 0644)

	// Test Valid JSON
	formatted, err := FormatJSON(validJSONFile)
	if err != nil {
		t.Errorf("FormatJSON failed for valid JSON: %v", err)
	}
	if !strings.Contains(formatted, `"hello": "world"`) {
		t.Errorf("Formatted JSON does not contain expected output: %s", formatted)
	}

	// Test Invalid JSON
	_, err = FormatJSON(invalidJSONFile)
	if err == nil {
		t.Errorf("Expected error for invalid JSON, got nil")
	}

	// Test Empty JSON
	_, err = FormatJSON(emptyJSONFile)
	if err == nil {
		t.Errorf("Expected error for empty JSON, got nil")
	}

	// Test Non-existent file
	_, err = FormatJSON(filepath.Join(tempDir, "nonexistent.json"))
	if err == nil {
		t.Errorf("Expected error for non-existent file, got nil")
	}

	// Test Invalid extension
	txtFile := filepath.Join(tempDir, "file.txt")
	os.WriteFile(txtFile, []byte(`{}`), 0644)
	_, err = FormatJSON(txtFile)
	if err == nil {
		t.Errorf("Expected error for non-JSON file extension, got nil")
	}
}
