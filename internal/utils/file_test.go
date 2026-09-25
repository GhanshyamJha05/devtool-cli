package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestOrganizeFolder(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "devtool-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir) // clean up

	// Create dummy files
	files := []string{
		"test1.jpg",
		"test2.png",
		"doc1.pdf",
		"doc2.txt",
		"video.mp4",
		"code.go",
		"archive.zip",
		"noextfile",
	}

	for _, f := range files {
		filePath := filepath.Join(tempDir, f)
		if err := os.WriteFile(filePath, []byte("dummy data"), 0644); err != nil {
			t.Fatalf("Failed to create dummy file %s: %v", f, err)
		}
	}

	// Run OrganizeFolder
	result, err := OrganizeFolder(tempDir, false)
	if err != nil {
		t.Fatalf("OrganizeFolder failed: %v", err)
	}

	// Verify the result
	if result.TotalFiles != 7 {
		t.Errorf("Expected 7 files to be moved, got %d", result.TotalFiles)
	}
	if result.Skipped != 1 {
		t.Errorf("Expected 1 file to be skipped, got %d", result.Skipped)
	}

	expectedCategories := map[string]int{
		"Images":    2,
		"Documents": 2,
		"Videos":    1,
		"Code":      1,
		"Archives":  1,
	}

	for category, expectedCount := range expectedCategories {
		if count := len(result.Moved[category]); count != expectedCount {
			t.Errorf("Expected %d files in %s, got %d", expectedCount, category, count)
		}
		
		// Verify folder exists
		categoryPath := filepath.Join(tempDir, category)
		info, err := os.Stat(categoryPath)
		if err != nil {
			t.Errorf("Expected category folder %s to exist, but got error: %v", category, err)
		} else if !info.IsDir() {
			t.Errorf("Expected %s to be a directory", categoryPath)
		}
	}
}
