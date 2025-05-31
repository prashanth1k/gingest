package ingester

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/prashanth1k/gingest/internal/types"
)

// ProcessRemoteRepo clones a Git repository and processes its files
func ProcessRemoteRepo(gitURL string, targetBranch string) (string, []types.FileInfo, error) {
	tempDir, filesData, _, err := ProcessRemoteRepoWithOptions(gitURL, targetBranch, 0) // 0 means no size limit
	return tempDir, filesData, err
}

// ProcessRemoteRepoWithOptions clones a Git repository and processes its files with size filtering
func ProcessRemoteRepoWithOptions(gitURL string, targetBranch string, maxFileSize int64) (string, []types.FileInfo, types.Stats, error) {
	return ProcessRemoteRepoWithPatterns(gitURL, targetBranch, maxFileSize, nil, nil)
}

// ProcessRemoteRepoWithPatterns clones a Git repository and processes its files with patterns
func ProcessRemoteRepoWithPatterns(gitURL string, targetBranch string, maxFileSize int64, includePatterns, excludePatterns []string) (string, []types.FileInfo, types.Stats, error) {
	// Create temporary directory for cloning
	tempDir, err := os.MkdirTemp("", "gingest-clone-*")
	if err != nil {
		return "", nil, types.Stats{}, fmt.Errorf("failed to create temp directory: %w", err)
	}

	// Defer cleanup
	defer func() {
		os.RemoveAll(tempDir)
	}()

	// Construct git clone command
	args := []string{"clone", "--depth", "1"}

	// Add branch-specific arguments if targetBranch is specified
	if targetBranch != "" {
		args = append(args, "-b", targetBranch, "--single-branch")
	}

	args = append(args, gitURL, tempDir)

	// Execute git clone
	cmd := exec.Command("git", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", nil, types.Stats{}, fmt.Errorf("git clone failed: %w\nOutput: %s", err, string(output))
	}

	// Process the cloned directory with size filtering and patterns
	filesData, stats, err := ProcessLocalDirectoryWithPatterns(tempDir, maxFileSize, includePatterns, excludePatterns)
	if err != nil {
		return "", nil, types.Stats{}, fmt.Errorf("failed to process cloned directory: %w", err)
	}

	// Update stats with Git-specific information
	stats.Source = gitURL
	stats.Branch = targetBranch

	return tempDir, filesData, stats, nil
}
