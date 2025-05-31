//go:build integration
// +build integration

package gingest

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestIntegration(t *testing.T) {
	// Create a temporary test directory
	testDir := t.TempDir()

	// Create diverse test files
	createTestFiles(t, testDir)

	// Build gingest CLI
	execName := "gingest_test"
	if runtime.GOOS == "windows" {
		execName = "gingest_test.exe"
	}
	buildCmd := exec.Command("go", "build", "-o", execName, "cmd/gingest/main.go")
	err := buildCmd.Run()
	if err != nil {
		t.Fatalf("Failed to build gingest: %v", err)
	}
	defer os.Remove(execName)

	// Test basic functionality
	outputFile := filepath.Join(testDir, "digest.md")
	execPath := "./" + execName
	if runtime.GOOS == "windows" {
		execPath = ".\\" + execName
	}
	cmd := exec.Command(execPath, "--source="+testDir, "--output="+outputFile)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to run gingest: %v\nOutput: %s", err, string(output))
	}

	// Verify output file exists
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Fatal("Output file was not created")
	}

	// Read and verify digest content
	digestContent, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read digest file: %v", err)
	}

	digestStr := string(digestContent)

	// Verify summary section
	if !strings.Contains(digestStr, "# Codebase Digest Summary") {
		t.Error("Summary header not found in digest")
	}

	if !strings.Contains(digestStr, "## Statistics") {
		t.Error("Statistics section not found in digest")
	}

	if !strings.Contains(digestStr, "## Directory Structure") {
		t.Error("Directory structure not found in digest")
	}

	// Verify file content is included
	if !strings.Contains(digestStr, "FILE: test.txt") {
		t.Error("test.txt file not found in digest")
	}

	if !strings.Contains(digestStr, "Hello, World!") {
		t.Error("test.txt content not found in digest")
	}

	if !strings.Contains(digestStr, "FILE: test.go") {
		t.Error("test.go file not found in digest")
	}

	if !strings.Contains(digestStr, "package main") {
		t.Error("test.go content not found in digest")
	}

	// Verify binary file handling
	if !strings.Contains(digestStr, "FILE: binary.bin") {
		t.Error("binary.bin file not found in digest")
	}

	if !strings.Contains(digestStr, "[Binary File]") {
		t.Error("Binary file marker not found in digest")
	}

	// Verify README prioritization (should appear before other files)
	readmeIndex := strings.Index(digestStr, "FILE: README.md")
	testTxtIndex := strings.Index(digestStr, "FILE: test.txt")
	if readmeIndex == -1 {
		t.Error("README.md not found in digest")
	}
	if testTxtIndex == -1 {
		t.Error("test.txt not found in digest")
	}
	if readmeIndex > testTxtIndex {
		t.Error("README.md should appear before other files")
	}

	// Verify Jupyter notebook processing
	if !strings.Contains(digestStr, "FILE: test.ipynb") {
		t.Error("test.ipynb file not found in digest")
	}

	if !strings.Contains(digestStr, "# Jupyter Notebook Content") {
		t.Error("Jupyter notebook content header not found")
	}

	if !strings.Contains(digestStr, "Integration test notebook") {
		t.Error("Jupyter notebook markdown content not found")
	}
}

func TestIntegrationWithPatterns(t *testing.T) {
	// Create a temporary test directory
	testDir := t.TempDir()

	// Create diverse test files
	createTestFiles(t, testDir)

	// Build gingest CLI
	execName := "gingest_test"
	if runtime.GOOS == "windows" {
		execName = "gingest_test.exe"
	}
	buildCmd := exec.Command("go", "build", "-o", execName, "cmd/gingest/main.go")
	err := buildCmd.Run()
	if err != nil {
		t.Fatalf("Failed to build gingest: %v", err)
	}
	defer os.Remove(execName)

	// Test with include patterns (only .go files)
	outputFile := filepath.Join(testDir, "go_only_digest.md")
	execPath := "./" + execName
	if runtime.GOOS == "windows" {
		execPath = ".\\" + execName
	}
	cmd := exec.Command(execPath, "--source="+testDir, "--output="+outputFile, "--include=*.go")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to run gingest with patterns: %v\nOutput: %s", err, string(output))
	}

	// Read and verify digest content
	digestContent, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read digest file: %v", err)
	}

	digestStr := string(digestContent)

	// Should include .go files
	if !strings.Contains(digestStr, "FILE: test.go") {
		t.Error("test.go file should be included with *.go pattern")
	}

	// Should not include .txt files
	if strings.Contains(digestStr, "FILE: test.txt") {
		t.Error("test.txt file should not be included with *.go pattern")
	}

	// Should not include .ipynb files
	if strings.Contains(digestStr, "FILE: test.ipynb") {
		t.Error("test.ipynb file should not be included with *.go pattern")
	}
}

func TestIntegrationConcurrency(t *testing.T) {
	// Create a temporary test directory with many files to test concurrency
	testDir := t.TempDir()

	// Create multiple files to test concurrent processing
	for i := 0; i < 20; i++ {
		fileName := filepath.Join(testDir, fmt.Sprintf("file_%d.txt", i))
		content := fmt.Sprintf("Content of file %d\nLine 2\nLine 3", i)
		err := os.WriteFile(fileName, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file %d: %v", i, err)
		}
	}

	// Build gingest CLI
	execName := "gingest_test"
	if runtime.GOOS == "windows" {
		execName = "gingest_test.exe"
	}
	buildCmd := exec.Command("go", "build", "-o", execName, "cmd/gingest/main.go")
	err := buildCmd.Run()
	if err != nil {
		t.Fatalf("Failed to build gingest: %v", err)
	}
	defer os.Remove(execName)

	// Test concurrent processing
	outputFile := filepath.Join(testDir, "concurrent_digest.md")
	execPath := "./" + execName
	if runtime.GOOS == "windows" {
		execPath = ".\\" + execName
	}
	cmd := exec.Command(execPath, "--source="+testDir, "--output="+outputFile)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to run gingest: %v\nOutput: %s", err, string(output))
	}

	// Verify all files were processed
	digestContent, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read digest file: %v", err)
	}

	digestStr := string(digestContent)

	// Check that all files are present
	for i := 0; i < 20; i++ {
		expectedFile := fmt.Sprintf("FILE: file_%d.txt", i)
		if !strings.Contains(digestStr, expectedFile) {
			t.Errorf("File %d not found in digest", i)
		}

		expectedContent := fmt.Sprintf("Content of file %d", i)
		if !strings.Contains(digestStr, expectedContent) {
			t.Errorf("Content of file %d not found in digest", i)
		}
	}
}

func createTestFiles(t *testing.T, testDir string) {
	// Create a text file
	textFile := filepath.Join(testDir, "test.txt")
	err := os.WriteFile(textFile, []byte("Hello, World!\nThis is a test file."), 0644)
	if err != nil {
		t.Fatalf("Failed to create test.txt: %v", err)
	}

	// Create a Go file
	goFile := filepath.Join(testDir, "test.go")
	goContent := `package main

import "fmt"

func main() {
	fmt.Println("Hello from Go!")
}
`
	err = os.WriteFile(goFile, []byte(goContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test.go: %v", err)
	}

	// Create a README file
	readmeFile := filepath.Join(testDir, "README.md")
	readmeContent := `# Test Project

This is a test project for gingest integration testing.

## Features

- Text files
- Go files
- Binary files
- Jupyter notebooks
`
	err = os.WriteFile(readmeFile, []byte(readmeContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create README.md: %v", err)
	}

	// Create a binary file (with null bytes)
	binaryFile := filepath.Join(testDir, "binary.bin")
	binaryContent := []byte{0x00, 0x01, 0x02, 0x03, 0xFF, 0xFE, 0xFD}
	err = os.WriteFile(binaryFile, binaryContent, 0644)
	if err != nil {
		t.Fatalf("Failed to create binary.bin: %v", err)
	}

	// Create a subdirectory with files
	subDir := filepath.Join(testDir, "subdir")
	err = os.Mkdir(subDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create subdirectory: %v", err)
	}

	subFile := filepath.Join(subDir, "nested.txt")
	err = os.WriteFile(subFile, []byte("Nested file content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create nested.txt: %v", err)
	}

	// Create a Jupyter notebook
	notebookFile := filepath.Join(testDir, "test.ipynb")
	notebookContent := `{
 "cells": [
  {
   "cell_type": "markdown",
   "metadata": {},
   "source": [
    "# Test Notebook\n",
    "\n",
    "Integration test notebook."
   ]
  },
  {
   "cell_type": "code",
   "execution_count": null,
   "metadata": {},
   "source": [
    "print('Integration test')"
   ]
  }
 ],
 "metadata": {
  "kernelspec": {
   "display_name": "Python 3",
   "language": "python",
   "name": "python3"
  }
 },
 "nbformat": 4,
 "nbformat_minor": 4
}`
	err = os.WriteFile(notebookFile, []byte(notebookContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test.ipynb: %v", err)
	}
}
