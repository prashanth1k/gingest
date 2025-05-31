package ingester

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sync"

	"github.com/prashanth1k/gingest/internal/notebookparser"
	"github.com/prashanth1k/gingest/internal/types"
	"github.com/prashanth1k/gingest/internal/utils"
)

// Package ingester handles processing of local directories and remote repositories

// ProcessLocalDirectory traverses a directory and returns FileInfo for all files
func ProcessLocalDirectory(rootDir string) ([]types.FileInfo, error) {
	filesData, _, err := ProcessLocalDirectoryWithOptions(rootDir, 0) // 0 means no size limit
	return filesData, err
}

// ProcessLocalDirectoryWithOptions traverses a directory with size filtering and returns stats
func ProcessLocalDirectoryWithOptions(rootDir string, maxFileSize int64) ([]types.FileInfo, types.Stats, error) {
	return ProcessLocalDirectoryWithPatterns(rootDir, maxFileSize, nil, nil)
}

// ProcessLocalDirectoryWithPatterns traverses a directory with filtering patterns
func ProcessLocalDirectoryWithPatterns(rootDir string, maxFileSize int64, includePatterns, excludePatterns []string) ([]types.FileInfo, types.Stats, error) {
	var allPaths []string  // Collect all paths for tree generation
	var filePaths []string // Collect file paths for concurrent processing
	stats := types.Stats{
		Source: rootDir,
	}

	// First pass: collect all valid file paths
	err := filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Calculate relative path from rootDir for pattern matching
		relPath, err := filepath.Rel(rootDir, path)
		if err != nil {
			relPath = path // fallback to original path
		}
		relPath = filepath.ToSlash(relPath)

		// Skip the root directory itself
		if relPath == "." {
			return nil
		}

		// Count directories
		if d.IsDir() {
			stats.NumDirsProcessed++

			// Check if directory should be excluded
			if !utils.ShouldIncludeFile(relPath, includePatterns, excludePatterns) {
				return filepath.SkipDir
			}
			return nil
		}

		// Check if file should be included based on patterns
		if !utils.ShouldIncludeFile(relPath, includePatterns, excludePatterns) {
			return nil // Skip this file
		}

		// Add to paths for tree generation
		allPaths = append(allPaths, relPath)

		// Get absolute path for processing
		absPath, err := filepath.Abs(path)
		if err != nil {
			return err
		}

		// Add to file paths for concurrent processing
		filePaths = append(filePaths, absPath)

		return nil
	})

	if err != nil {
		return nil, types.Stats{}, err
	}

	// Second pass: process files concurrently
	filesData := make([]types.FileInfo, len(filePaths))
	fileInfoChan := make(chan struct {
		index int
		info  types.FileInfo
	}, len(filePaths))

	var wg sync.WaitGroup
	var statsMutex sync.Mutex

	// Process each file concurrently
	for i, absPath := range filePaths {
		wg.Add(1)
		go func(index int, filePath string) {
			defer wg.Done()

			// Calculate relative path
			relPath, err := filepath.Rel(rootDir, filePath)
			if err != nil {
				relPath = filePath
			}
			relPath = filepath.ToSlash(relPath)

			// Get file info to check size
			fileInfo, err := os.Stat(filePath)
			if err != nil {
				fileInfoStruct := types.FileInfo{
					RelativePath: relPath,
					AbsolutePath: filePath,
					Content:      "",
					IsBinary:     false,
					Error:        err,
				}
				fileInfoChan <- struct {
					index int
					info  types.FileInfo
				}{index, fileInfoStruct}

				statsMutex.Lock()
				stats.NumFilesProcessed++
				statsMutex.Unlock()
				return
			}

			var content string
			var readErr error
			var isBinary bool

			// Check file size if maxFileSize is specified
			if maxFileSize > 0 && fileInfo.Size() > maxFileSize {
				sizeMB := float64(fileInfo.Size()) / (1024 * 1024)
				content = fmt.Sprintf("[File content skipped: Exceeds max size (%.1f MB > %.1f MB)]",
					sizeMB, float64(maxFileSize)/(1024*1024))

				statsMutex.Lock()
				stats.NumSkippedFiles++
				statsMutex.Unlock()
			} else {
				// Check if file is a Jupyter notebook first
				if utils.IsJupyterNotebook(filePath) {
					content, readErr = notebookparser.ParseNotebook(filePath)
					if readErr == nil {
						statsMutex.Lock()
						stats.TotalContentBytes += int64(len(content))
						statsMutex.Unlock()
						isBinary = false // Notebooks are treated as text
					}
				} else {
					// Check if file is binary before reading full content
					isBinary, err := utils.IsBinaryFile(filePath)
					if err != nil {
						readErr = err
					} else if isBinary {
						content = "[Binary File]"
						statsMutex.Lock()
						stats.NumBinaryFiles++
						statsMutex.Unlock()
					} else {
						// Read file content for text files
						content, readErr = utils.ReadFileContent(filePath)
						statsMutex.Lock()
						stats.TotalContentBytes += int64(len(content))
						statsMutex.Unlock()
					}
				}
			}

			fileInfoStruct := types.FileInfo{
				RelativePath: relPath,
				AbsolutePath: filePath,
				Content:      content,
				IsBinary:     isBinary,
				Error:        readErr,
			}

			fileInfoChan <- struct {
				index int
				info  types.FileInfo
			}{index, fileInfoStruct}

			statsMutex.Lock()
			stats.NumFilesProcessed++
			statsMutex.Unlock()
		}(i, absPath)
	}

	// Collect results from channel
	go func() {
		wg.Wait()
		close(fileInfoChan)
	}()

	// Collect all file info structs in order
	for result := range fileInfoChan {
		filesData[result.index] = result.info
	}

	// Store all paths in stats for tree generation
	stats.AllPaths = allPaths

	return filesData, stats, nil
}
