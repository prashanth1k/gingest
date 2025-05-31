package notebookparser

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseNotebook(t *testing.T) {
	// Create a temporary notebook file
	tempDir := t.TempDir()
	notebookPath := filepath.Join(tempDir, "test.ipynb")

	// Sample notebook JSON
	notebookJSON := `{
 "cells": [
  {
   "cell_type": "markdown",
   "metadata": {},
   "source": [
    "# Test Notebook\n",
    "\n",
    "This is a test."
   ]
  },
  {
   "cell_type": "code",
   "execution_count": null,
   "metadata": {},
   "source": [
    "print('Hello, World!')\n",
    "x = 42"
   ]
  },
  {
   "cell_type": "raw",
   "metadata": {},
   "source": [
    "Raw cell content"
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

	// Write the notebook file
	err := os.WriteFile(notebookPath, []byte(notebookJSON), 0644)
	if err != nil {
		t.Fatalf("Failed to create test notebook: %v", err)
	}

	// Parse the notebook
	content, err := ParseNotebook(notebookPath)
	if err != nil {
		t.Fatalf("Failed to parse notebook: %v", err)
	}

	// Verify the content
	if !strings.Contains(content, "# Jupyter Notebook Content") {
		t.Error("Expected header not found in parsed content")
	}

	if !strings.Contains(content, "## Cell 1 (markdown)") {
		t.Error("Expected markdown cell header not found")
	}

	if !strings.Contains(content, "# Test Notebook") {
		t.Error("Expected markdown content not found")
	}

	if !strings.Contains(content, "## Cell 2 (code)") {
		t.Error("Expected code cell header not found")
	}

	if !strings.Contains(content, "print('Hello, World!')") {
		t.Error("Expected code content not found")
	}

	if !strings.Contains(content, "## Cell 3 (raw)") {
		t.Error("Expected raw cell header not found")
	}

	if !strings.Contains(content, "Raw cell content") {
		t.Error("Expected raw content not found")
	}
}

func TestParseNotebook_InvalidJSON(t *testing.T) {
	// Create a temporary file with invalid JSON
	tempDir := t.TempDir()
	notebookPath := filepath.Join(tempDir, "invalid.ipynb")

	invalidJSON := `{"cells": [invalid json}`

	err := os.WriteFile(notebookPath, []byte(invalidJSON), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Try to parse the invalid notebook
	_, err = ParseNotebook(notebookPath)
	if err == nil {
		t.Error("Expected error for invalid JSON, but got none")
	}

	if !strings.Contains(err.Error(), "failed to parse notebook JSON") {
		t.Errorf("Expected JSON parse error, got: %v", err)
	}
}

func TestParseNotebook_NonexistentFile(t *testing.T) {
	// Try to parse a file that doesn't exist
	_, err := ParseNotebook("nonexistent.ipynb")
	if err == nil {
		t.Error("Expected error for nonexistent file, but got none")
	}

	if !strings.Contains(err.Error(), "failed to read notebook file") {
		t.Errorf("Expected file read error, got: %v", err)
	}
}
