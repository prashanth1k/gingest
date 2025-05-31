// Package gingest provides functionality to convert codebases into LLM-friendly text digests.
// It supports both local directories and remote Git repositories.
package gingest

import (
	"fmt"
	"os"

	"github.com/prashanth1k/gingest/internal/ingester"
	"github.com/prashanth1k/gingest/internal/types"
	"github.com/prashanth1k/gingest/internal/utils"
)

// Config holds configuration options for processing codebases
type Config struct {
	Source       string // Local directory path or Git repository URL
	OutputFile   string // Output file path for the digest
	TargetBranch string // Target branch for Git repositories (optional)
	MaxFileSize  int64  // Maximum file size in bytes (0 = no limit)
}

// ProcessAndWriteDigest processes a codebase and writes the digest to a file
func ProcessAndWriteDigest(config Config) error {
	filesData, stats, err := ProcessCodebaseWithStats(config)
	if err != nil {
		return fmt.Errorf("failed to process codebase: %w", err)
	}

	err = ingester.WriteDigest(config.OutputFile, filesData, stats)
	if err != nil {
		return fmt.Errorf("failed to write digest: %w", err)
	}

	return nil
}

// ProcessCodebase processes a codebase and returns file information
func ProcessCodebase(config Config) ([]types.FileInfo, error) {
	filesData, _, err := ProcessCodebaseWithStats(config)
	return filesData, err
}

// ProcessCodebaseWithStats processes a codebase and returns file information with statistics
func ProcessCodebaseWithStats(config Config) ([]types.FileInfo, types.Stats, error) {
	var filesData []types.FileInfo
	var stats types.Stats
	var err error

	// Check if source is a Git URL
	if utils.IsGitURL(config.Source) {
		_, filesData, stats, err = ingester.ProcessRemoteRepoWithOptions(
			config.Source,
			config.TargetBranch,
			config.MaxFileSize,
		)
		if err != nil {
			return nil, types.Stats{}, fmt.Errorf("failed to process remote repository: %w", err)
		}
	} else if info, err := os.Stat(config.Source); err == nil && info.IsDir() {
		filesData, stats, err = ingester.ProcessLocalDirectoryWithOptions(
			config.Source,
			config.MaxFileSize,
		)
		if err != nil {
			return nil, types.Stats{}, fmt.Errorf("failed to process local directory: %w", err)
		}
	} else {
		return nil, types.Stats{}, fmt.Errorf("source must be a valid local directory or Git URL: %s", config.Source)
	}

	return filesData, stats, nil
}

// FileInfo represents information about a processed file (exported version)
type FileInfo = types.FileInfo

// Stats represents processing statistics (exported version)
type Stats = types.Stats
