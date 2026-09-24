package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// CleanResult holds a summary of the clean operation for display.
type CleanResult struct {
	TotalFiles int
	Skipped    int
	Moved      map[string][]string // category → list of filenames
}

// OrganizeFolder reads a directory and moves files into categorized subfolders.
// Returns a structured result for the caller to display.
func OrganizeFolder(targetDir string) (*CleanResult, error) {
	// Step 1: Validate that the path exists and is a directory
	info, err := os.Stat(targetDir)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("directory not found: '%s'", targetDir)
	}
	if err != nil {
		return nil, fmt.Errorf("cannot access path '%s': %w", targetDir, err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("'%s' is not a directory", targetDir)
	}

	// Step 2: Read entries
	entries, err := os.ReadDir(targetDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory '%s': %w", targetDir, err)
	}

	result := &CleanResult{
		Moved: make(map[string][]string),
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		if entry.Type()&os.ModeSymlink != 0 {
			result.Skipped++
			continue
		}

		fileName := entry.Name()
		ext := strings.ToLower(filepath.Ext(fileName))

		// Skip files without extensions
		if ext == "" {
			result.Skipped++
			continue
		}

		result.TotalFiles++
		category := getCategoryForExtension(ext)
		destFolder := filepath.Join(targetDir, category)

		// Create category folder if needed
		if err := os.MkdirAll(destFolder, 0755); err != nil {
			return nil, fmt.Errorf("failed to create folder '%s': %w", destFolder, err)
		}

		oldPath := filepath.Join(targetDir, fileName)
		newPath := filepath.Join(destFolder, fileName)

		// Collision handling
		baseName := strings.TrimSuffix(fileName, ext)
		counter := 1
		for {
			if _, err := os.Stat(newPath); os.IsNotExist(err) {
				break
			}
			newFileName := fmt.Sprintf("%s (%d)%s", baseName, counter, ext)
			newPath = filepath.Join(destFolder, newFileName)
			counter++
		}

		if err := os.Rename(oldPath, newPath); err != nil {
			return nil, fmt.Errorf("failed to move '%s': %w", fileName, err)
		}

		result.Moved[category] = append(result.Moved[category], fileName)
	}

	return result, nil
}

// getCategoryForExtension maps file extensions to human-readable folder names.
func getCategoryForExtension(ext string) string {
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".svg", ".webp", ".ico", ".bmp":
		return "Images"
	case ".pdf", ".doc", ".docx", ".txt", ".md", ".csv", ".xlsx", ".pptx":
		return "Documents"
	case ".mp4", ".mkv", ".avi", ".mov", ".wmv", ".flv":
		return "Videos"
	case ".mp3", ".wav", ".flac", ".aac", ".ogg":
		return "Audio"
	case ".zip", ".tar", ".gz", ".rar", ".7z":
		return "Archives"
	case ".go", ".js", ".ts", ".py", ".java", ".json", ".html", ".css", ".cpp", ".c", ".rs", ".rb":
		return "Code"
	case ".exe", ".msi", ".dmg", ".deb", ".rpm":
		return "Executables"
	default:
		return "Others"
	}
}
