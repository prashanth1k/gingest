package utils

import (
	"os/exec"
	"strings"
)

// IsGitURL checks if a source string is a Git URL
func IsGitURL(source string) bool {
	// Check for SSH Git URLs (git@host:repo.git)
	if strings.HasPrefix(source, "git@") {
		return true
	}

	// Check for .git suffix
	if strings.HasSuffix(source, ".git") {
		return true
	}

	// Check for common Git hosting services
	if strings.Contains(source, "github.com") ||
		strings.Contains(source, "gitlab.com") ||
		strings.Contains(source, "bitbucket.org") {
		return true
	}

	return false
}

// IsGitAvailable checks if git command is available on PATH
func IsGitAvailable() bool {
	_, err := exec.LookPath("git")
	return err == nil
}
