package ingester

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/prashanth1k/gingest/internal/types"
	"github.com/prashanth1k/gingest/internal/utils"
)

const (
	FILE_SEPARATOR_START = "================================================"
	FILE_SEPARATOR_END   = "================================================"
)

// WriteDigest writes the collected file data to the output file with summary
func WriteDigest(outputFilePath string, filesData []types.FileInfo, stats types.Stats) error {
	// Partition files into README and other files
	var readmeFiles []types.FileInfo
	var otherFiles []types.FileInfo

	for _, fileInfo := range filesData {
		// Skip files with errors
		if fileInfo.Error != nil {
			continue
		}

		fileName := filepath.Base(fileInfo.RelativePath)
		if utils.IsReadmeFile(fileName) {
			readmeFiles = append(readmeFiles, fileInfo)
		} else {
			otherFiles = append(otherFiles, fileInfo)
		}
	}

	// Sort both groups by RelativePath
	sort.Slice(readmeFiles, func(i, j int) bool {
		return readmeFiles[i].RelativePath < readmeFiles[j].RelativePath
	})
	sort.Slice(otherFiles, func(i, j int) bool {
		return otherFiles[i].RelativePath < otherFiles[j].RelativePath
	})

	// Create output file
	file, err := os.Create(outputFilePath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer file.Close()

	// Write summary at the beginning
	summary := utils.GenerateSummaryString(stats)
	_, err = file.WriteString(summary)
	if err != nil {
		return fmt.Errorf("failed to write summary: %w", err)
	}

	// Write directory tree if we have paths
	if len(stats.AllPaths) > 0 {
		rootName := filepath.Base(stats.Source)
		if rootName == "." || rootName == "" {
			rootName = "project"
		}
		tree := utils.GenerateTreeString(stats.AllPaths, rootName, filesData)
		_, err = file.WriteString("## Directory Structure\n\n```\n")
		if err != nil {
			return fmt.Errorf("failed to write tree header: %w", err)
		}
		_, err = file.WriteString(tree)
		if err != nil {
			return fmt.Errorf("failed to write tree: %w", err)
		}
		_, err = file.WriteString("```\n\n---\n\n")
		if err != nil {
			return fmt.Errorf("failed to write tree footer: %w", err)
		}
	}

	// Write README files first, then other files
	allFiles := append(readmeFiles, otherFiles...)

	// Write each file's content
	for _, fileInfo := range allFiles {
		// Write file separator and header
		_, err := fmt.Fprintf(file, "%s\n", FILE_SEPARATOR_START)
		if err != nil {
			return fmt.Errorf("failed to write separator: %w", err)
		}

		_, err = fmt.Fprintf(file, "FILE: %s\n", fileInfo.RelativePath)
		if err != nil {
			return fmt.Errorf("failed to write file header: %w", err)
		}

		_, err = fmt.Fprintf(file, "%s\n", FILE_SEPARATOR_END)
		if err != nil {
			return fmt.Errorf("failed to write separator: %w", err)
		}

		// Write file content
		_, err = fmt.Fprintf(file, "%s", fileInfo.Content)
		if err != nil {
			return fmt.Errorf("failed to write file content: %w", err)
		}

		// Write two newlines after content
		_, err = fmt.Fprintf(file, "\n\n")
		if err != nil {
			return fmt.Errorf("failed to write newlines: %w", err)
		}
	}

	return nil
}
