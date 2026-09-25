package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// FormatJSON reads a JSON file, validates it, and writes a pretty-printed version to the provided writer.
func FormatJSON(filePath string, out io.Writer) error {
	// Step 1: Validate the file path before doing any I/O
	if err := validateJSONFile(filePath); err != nil {
		return err
	}

	// Step 2: Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file '%s': %w", filePath, err)
	}
	defer file.Close()

	// Step 3: Check for empty file
	info, err := file.Stat()
	if err == nil && info.Size() == 0 {
		return fmt.Errorf("file '%s' is empty — nothing to format", filePath)
	}

	// Step 4: Unmarshal to validate JSON structure
	var parsedJSON interface{}
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&parsedJSON); err != nil {
		if err == io.EOF {
			return fmt.Errorf("file '%s' is empty — nothing to format", filePath)
		}
		return fmt.Errorf("invalid JSON in '%s': %w", filePath, err)
	}

	// Step 5: Encode with indentation directly to writer
	encoder := json.NewEncoder(out)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(parsedJSON); err != nil {
		return fmt.Errorf("failed to format JSON: %w", err)
	}

	return nil
}

// validateJSONFile checks that the file exists and has a .json extension.
func validateJSONFile(filePath string) error {
	// Check if file exists
	info, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return fmt.Errorf("file not found: '%s'", filePath)
	}
	if err != nil {
		return fmt.Errorf("cannot access file '%s': %w", filePath, err)
	}

	// Don't allow directories
	if info.IsDir() {
		return fmt.Errorf("'%s' is a directory, not a file", filePath)
	}

	// Warn if not a .json extension
	ext := strings.ToLower(filepath.Ext(filePath))
	if ext != ".json" {
		return fmt.Errorf("expected a .json file, got '%s'", ext)
	}

	return nil
}
