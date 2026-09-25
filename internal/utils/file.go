package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
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

	var mu sync.Mutex
	var wg sync.WaitGroup
	var moveErr error

	type job struct {
		fileName string
		ext      string
	}

	jobs := make(chan job, len(entries))

	numWorkers := 5
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := range jobs {
				mu.Lock()
				hasErr := moveErr != nil
				mu.Unlock()
				if hasErr {
					continue
				}

				category := getCategoryForExtension(j.ext)
				destFolder := filepath.Join(targetDir, category)

				if err := os.MkdirAll(destFolder, 0755); err != nil {
					mu.Lock()
					if moveErr == nil {
						moveErr = fmt.Errorf("failed to create folder '%s': %w", destFolder, err)
					}
					mu.Unlock()
					continue
				}

				oldPath := filepath.Join(targetDir, j.fileName)
				newPath := filepath.Join(destFolder, j.fileName)

				baseName := strings.TrimSuffix(j.fileName, j.ext)
				counter := 1
				for {
					if _, err := os.Stat(newPath); os.IsNotExist(err) {
						break
					}
					newFileName := fmt.Sprintf("%s (%d)%s", baseName, counter, j.ext)
					newPath = filepath.Join(destFolder, newFileName)
					counter++
				}

				if err := os.Rename(oldPath, newPath); err != nil {
					mu.Lock()
					if moveErr == nil {
						moveErr = fmt.Errorf("failed to move '%s': %w", j.fileName, err)
					}
					mu.Unlock()
					continue
				}

				mu.Lock()
				result.Moved[category] = append(result.Moved[category], j.fileName)
				result.TotalFiles++
				mu.Unlock()
			}
		}()
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

		if ext == "" {
			result.Skipped++
			continue
		}

		jobs <- job{fileName: fileName, ext: ext}
	}
	close(jobs)
	wg.Wait()

	if moveErr != nil {
		return result, moveErr
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
