package gingest

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestProcessCodebase_LocalDirectory(t *testing.T) {
	// Create a temporary directory with test files
	tempDir, err := os.MkdirTemp("", "gingest-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test files
	testFiles := map[string]string{
		"test1.txt": "Hello, World!",
		"test2.go":  "package main\n\nfunc main() {\n\tprintln(\"Hello\")\n}",
		"README.md": "# Test Project\n\nThis is a test.",
	}

	for filename, content := range testFiles {
		filePath := filepath.Join(tempDir, filename)
		err := os.WriteFile(filePath, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file %s: %v", filename, err)
		}
	}

	// Test processing
	config := Config{
		Source:      tempDir,
		MaxFileSize: 1024 * 1024, // 1MB
	}

	filesData, err := ProcessCodebase(config)
	if err != nil {
		t.Fatalf("ProcessCodebase failed: %v", err)
	}

	if len(filesData) != len(testFiles) {
		t.Errorf("Expected %d files, got %d", len(testFiles), len(filesData))
	}

	// Check that all test files are present
	foundFiles := make(map[string]bool)
	for _, fileInfo := range filesData {
		foundFiles[filepath.Base(fileInfo.RelativePath)] = true

		if fileInfo.Error != nil {
			t.Errorf("File %s has error: %v", fileInfo.RelativePath, fileInfo.Error)
		}

		expectedContent := testFiles[filepath.Base(fileInfo.RelativePath)]
		if fileInfo.Content != expectedContent {
			t.Errorf("File %s content mismatch. Expected: %q, Got: %q",
				fileInfo.RelativePath, expectedContent, fileInfo.Content)
		}
	}

	for filename := range testFiles {
		if !foundFiles[filename] {
			t.Errorf("Expected file %s not found in results", filename)
		}
	}
}

func TestProcessCodebase_MaxFileSize(t *testing.T) {
	// Create a temporary directory with a large file
	tempDir, err := os.MkdirTemp("", "gingest-test-large-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a large file (1KB)
	largeContent := make([]byte, 1024)
	for i := range largeContent {
		largeContent[i] = 'A'
	}

	largeFilePath := filepath.Join(tempDir, "large.txt")
	err = os.WriteFile(largeFilePath, largeContent, 0644)
	if err != nil {
		t.Fatalf("Failed to create large test file: %v", err)
	}

	// Test with size limit smaller than the file
	config := Config{
		Source:      tempDir,
		MaxFileSize: 512, // 512 bytes limit
	}

	filesData, err := ProcessCodebase(config)
	if err != nil {
		t.Fatalf("ProcessCodebase failed: %v", err)
	}

	if len(filesData) != 1 {
		t.Fatalf("Expected 1 file, got %d", len(filesData))
	}

	fileInfo := filesData[0]
	if fileInfo.Error != nil {
		t.Errorf("File has error: %v", fileInfo.Error)
	}

	// Check that content was skipped due to size
	if !strings.Contains(fileInfo.Content, "File content skipped: Exceeds max size") {
		t.Errorf("Expected size skip message, got: %s", fileInfo.Content)
	}
}

func TestProcessAndWriteDigest(t *testing.T) {
	// Create a temporary directory with test files
	tempDir, err := os.MkdirTemp("", "gingest-test-write-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test file
	testFile := filepath.Join(tempDir, "test.txt")
	err = os.WriteFile(testFile, []byte("Test content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Create output file path
	outputFile := filepath.Join(tempDir, "digest.md")

	// Test processing and writing
	config := Config{
		Source:     tempDir,
		OutputFile: outputFile,
	}

	err = ProcessAndWriteDigest(config)
	if err != nil {
		t.Fatalf("ProcessAndWriteDigest failed: %v", err)
	}

	// Check that output file was created
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Errorf("Output file was not created: %s", outputFile)
	}

	// Read and verify output content
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "FILE: test.txt") {
		t.Errorf("Output file doesn't contain expected file header")
	}

	if !strings.Contains(contentStr, "Test content") {
		t.Errorf("Output file doesn't contain expected file content")
	}
}
