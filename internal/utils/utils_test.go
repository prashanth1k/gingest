package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadFileContent(t *testing.T) {
	// Create a temporary file
	tempDir, err := os.MkdirTemp("", "utils-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	testContent := "Hello, World!\nThis is a test file."
	testFile := filepath.Join(tempDir, "test.txt")

	err = os.WriteFile(testFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Test reading the file
	content, err := ReadFileContent(testFile)
	if err != nil {
		t.Fatalf("ReadFileContent failed: %v", err)
	}

	if content != testContent {
		t.Errorf("Content mismatch. Expected: %q, Got: %q", testContent, content)
	}
}

func TestReadFileContent_NonExistentFile(t *testing.T) {
	// Test reading a non-existent file
	_, err := ReadFileContent("/non/existent/file.txt")
	if err == nil {
		t.Error("Expected error for non-existent file, got nil")
	}
}

func TestIsGitURL(t *testing.T) {
	testCases := []struct {
		url      string
		expected bool
		desc     string
	}{
		{"https://github.com/user/repo.git", true, "GitHub HTTPS URL"},
		{"http://gitlab.com/user/repo.git", true, "GitLab HTTP URL"},
		{"https://bitbucket.org/user/repo", true, "Bitbucket URL"},
		{"git@github.com:user/repo.git", true, "SSH Git URL"},
		{"https://example.com/repo.git", true, "Generic Git URL with .git suffix"},
		{"https://github.com/user/repo", true, "GitHub URL without .git"},
		{"/local/path", false, "Local path"},
		{"./relative/path", false, "Relative path"},
		{"https://example.com/page", false, "Regular HTTPS URL"},
		{"ftp://example.com/file", false, "FTP URL"},
		{"", false, "Empty string"},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			result := IsGitURL(tc.url)
			if result != tc.expected {
				t.Errorf("IsGitURL(%q) = %v, expected %v", tc.url, result, tc.expected)
			}
		})
	}
}
