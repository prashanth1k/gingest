package utils

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/prashanth1k/gingest/internal/types"
)

// Package utils provides utility functions for file operations and processing

// ReadFileContent reads the content of a file and returns it as a string
func ReadFileContent(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

// IsBinaryFile checks if a file is binary by reading the first 1024 bytes
// and looking for null bytes (simple heuristic)
func IsBinaryFile(filePath string) (bool, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return false, err
	}
	defer file.Close()

	// Read first 1024 bytes
	chunk := make([]byte, 1024)
	n, err := file.Read(chunk)
	if err != nil && n == 0 {
		return false, err
	}

	// Check for null bytes
	return bytes.Contains(chunk[:n], []byte{0}), nil
}

// IsReadmeFile checks if a filename is a README file (case-insensitive)
func IsReadmeFile(fileName string) bool {
	lowerName := strings.ToLower(fileName)
	return strings.HasPrefix(lowerName, "readme")
}

// GenerateSummaryString formats processing statistics into a readable summary
func GenerateSummaryString(stats types.Stats) string {
	var summary strings.Builder

	summary.WriteString("# Codebase Digest Summary\n\n")
	summary.WriteString(fmt.Sprintf("**Source:** %s\n", stats.Source))

	if stats.Branch != "" {
		summary.WriteString(fmt.Sprintf("**Branch:** %s\n", stats.Branch))
	}

	summary.WriteString(fmt.Sprintf("**Generated:** %s\n\n", time.Now().Format("2006-01-02 15:04:05")))

	summary.WriteString("## Statistics\n\n")
	summary.WriteString(fmt.Sprintf("- **Total Files:** %d\n", stats.NumFilesProcessed))
	summary.WriteString(fmt.Sprintf("- **Directories:** %d\n", stats.NumDirsProcessed))
	summary.WriteString(fmt.Sprintf("- **Binary Files:** %d\n", stats.NumBinaryFiles))
	summary.WriteString(fmt.Sprintf("- **Skipped Files:** %d\n", stats.NumSkippedFiles))
	summary.WriteString(fmt.Sprintf("- **Total Content Size:** %.2f KB\n\n", float64(stats.TotalContentBytes)/1024))

	summary.WriteString("---\n\n")

	return summary.String()
}

// ParsePatterns parses a comma-separated string of patterns into a slice
func ParsePatterns(patternsString string) []string {
	if patternsString == "" {
		return nil
	}

	patterns := strings.Split(patternsString, ",")
	var result []string

	for _, pattern := range patterns {
		trimmed := strings.TrimSpace(pattern)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	return result
}

// ShouldIncludeFile determines if a file should be included based on include/exclude patterns
func ShouldIncludeFile(relativePath string, includePatterns, excludePatterns []string) bool {
	fileName := filepath.Base(relativePath)

	// Check exclude patterns first
	for _, pattern := range excludePatterns {
		// Check against full path
		if matched, _ := filepath.Match(pattern, relativePath); matched {
			// Check if include patterns override the exclusion
			if len(includePatterns) > 0 {
				for _, includePattern := range includePatterns {
					if matched, _ := filepath.Match(includePattern, relativePath); matched {
						return true // Include pattern overrides exclude
					}
					if matched, _ := filepath.Match(includePattern, fileName); matched {
						return true // Include pattern overrides exclude
					}
				}
			}
			return false // Excluded and no include override
		}

		// Check against filename only
		if matched, _ := filepath.Match(pattern, fileName); matched {
			// Check if include patterns override the exclusion
			if len(includePatterns) > 0 {
				for _, includePattern := range includePatterns {
					if matched, _ := filepath.Match(includePattern, relativePath); matched {
						return true // Include pattern overrides exclude
					}
					if matched, _ := filepath.Match(includePattern, fileName); matched {
						return true // Include pattern overrides exclude
					}
				}
			}
			return false // Excluded and no include override
		}

		// Also check if any parent directory matches exclude pattern
		dir := filepath.Dir(relativePath)
		for dir != "." && dir != "/" {
			if matched, _ := filepath.Match(pattern, filepath.Base(dir)); matched {
				return false
			}
			dir = filepath.Dir(dir)
		}
	}

	// If include patterns are specified, file must match at least one
	if len(includePatterns) > 0 {
		for _, pattern := range includePatterns {
			// Check against full path
			if matched, _ := filepath.Match(pattern, relativePath); matched {
				return true
			}
			// Check against filename only
			if matched, _ := filepath.Match(pattern, fileName); matched {
				return true
			}
		}
		return false // Doesn't match any include pattern
	}

	return true // No exclusion and no include restriction
}

// GenerateTreeString creates a tree representation of file and directory paths
func GenerateTreeString(paths []string, rootName string, filesData []types.FileInfo) string {
	if len(paths) == 0 {
		return ""
	}

	var tree strings.Builder
	tree.WriteString(fmt.Sprintf("%s/\n", rootName))

	// Sort paths to ensure consistent tree structure
	sort.Strings(paths)

	// Create a map for quick lookup of file info
	fileInfoMap := make(map[string]types.FileInfo)
	for _, fileInfo := range filesData {
		fileInfoMap[fileInfo.RelativePath] = fileInfo
	}

	// Track which directories we've seen at each level
	dirTracker := make(map[string]bool)

	for i, path := range paths {
		parts := strings.Split(path, "/")

		// Build directory structure first
		currentPath := ""
		for level := 0; level < len(parts)-1; level++ {
			if level == 0 {
				currentPath = parts[level]
			} else {
				currentPath = currentPath + "/" + parts[level]
			}

			if !dirTracker[currentPath] {
				dirTracker[currentPath] = true

				// Determine prefix
				prefix := strings.Repeat("│   ", level)

				// Check if this is the last directory at this level
				isLastDir := true
				for j := i + 1; j < len(paths); j++ {
					otherParts := strings.Split(paths[j], "/")
					if len(otherParts) > level+1 {
						otherPath := strings.Join(otherParts[:level+1], "/")
						if strings.HasPrefix(otherPath, strings.Join(parts[:level+1], "/")) {
							isLastDir = false
							break
						}
					}
				}

				if isLastDir {
					tree.WriteString(prefix + "└── " + parts[level] + "/\n")
				} else {
					tree.WriteString(prefix + "├── " + parts[level] + "/\n")
				}
			}
		}

		// Now add the file
		fileName := parts[len(parts)-1]
		prefix := strings.Repeat("│   ", len(parts)-1)

		// Check if this is the last file in its directory
		isLastFile := true
		currentDir := strings.Join(parts[:len(parts)-1], "/")
		for j := i + 1; j < len(paths); j++ {
			otherParts := strings.Split(paths[j], "/")
			otherDir := strings.Join(otherParts[:len(otherParts)-1], "/")
			if otherDir == currentDir {
				isLastFile = false
				break
			}
		}

		// Add file with appropriate suffix
		suffix := ""
		if fileInfo, exists := fileInfoMap[path]; exists {
			if fileInfo.IsBinary {
				suffix = " (Binary)"
			} else if strings.Contains(fileInfo.Content, "[File content skipped:") {
				suffix = " (Skipped - Too Large)"
			}
		}

		if isLastFile {
			tree.WriteString(prefix + "└── " + fileName + suffix + "\n")
		} else {
			tree.WriteString(prefix + "├── " + fileName + suffix + "\n")
		}
	}

	return tree.String()
}

// IsJupyterNotebook checks if a file is a Jupyter notebook by extension
func IsJupyterNotebook(filePath string) bool {
	return strings.HasSuffix(strings.ToLower(filePath), ".ipynb")
}
