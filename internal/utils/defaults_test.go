package utils

import (
	"testing"
)

func TestGetDefaultExcludePatterns(t *testing.T) {
	patterns := GetDefaultExcludePatterns()

	// Test that we have a comprehensive list
	if len(patterns) < 50 {
		t.Errorf("Expected at least 50 default exclusion patterns, got %d", len(patterns))
	}

	// Test for key exclusions that should be present
	expectedPatterns := []string{
		// Version control
		".git",
		".svn",

		// Dependencies - JavaScript/Node.js
		"node_modules",
		".npm",

		// Dependencies - Python
		".venv",
		"venv",
		"__pycache__",

		// Dependencies - Go
		"vendor",

		// Dependencies - Java
		"target",
		".gradle",

		// IDE files
		".vscode",
		".idea",

		// OS files
		".DS_Store",
		"Thumbs.db",

		// Binary files
		"*.exe",
		"*.dll",
		"*.so",

		// Media files
		"*.jpg",
		"*.mp4",
		"*.mp3",

		// Lock files
		"package-lock.json",
		"yarn.lock",
		"go.sum",
	}

	patternMap := make(map[string]bool)
	for _, pattern := range patterns {
		patternMap[pattern] = true
	}

	for _, expected := range expectedPatterns {
		if !patternMap[expected] {
			t.Errorf("Expected pattern '%s' not found in default exclusions", expected)
		}
	}
}

func TestGetDefaultDirectoryExclusions(t *testing.T) {
	patterns := GetDefaultDirectoryExclusions()

	// Should contain common dependency directories
	expectedDirs := []string{
		"node_modules",
		".venv",
		"venv",
		"vendor",
		"target",
		".git",
	}

	patternMap := make(map[string]bool)
	for _, pattern := range patterns {
		patternMap[pattern] = true
	}

	for _, expected := range expectedDirs {
		if !patternMap[expected] {
			t.Errorf("Expected directory pattern '%s' not found in directory exclusions", expected)
		}
	}
}

func TestGetDefaultFileExclusions(t *testing.T) {
	patterns := GetDefaultFileExclusions()

	// Should contain common file patterns
	expectedFiles := []string{
		"*.log",
		"*.tmp",
		"*.exe",
		"*.jpg",
		"*.mp4",
		"package-lock.json",
		"yarn.lock",
	}

	patternMap := make(map[string]bool)
	for _, pattern := range patterns {
		patternMap[pattern] = true
	}

	for _, expected := range expectedFiles {
		if !patternMap[expected] {
			t.Errorf("Expected file pattern '%s' not found in file exclusions", expected)
		}
	}
}

func TestAddCustomExclusion(t *testing.T) {
	// Get initial count
	initialPatterns := GetDefaultExcludePatterns()
	initialCount := len(initialPatterns)

	// Add a custom pattern
	customPattern := "*.custom"
	AddCustomExclusion(customPattern)

	// Check that it was added
	updatedPatterns := GetDefaultExcludePatterns()
	if len(updatedPatterns) != initialCount+1 {
		t.Errorf("Expected %d patterns after adding custom, got %d", initialCount+1, len(updatedPatterns))
	}

	// Check that the custom pattern is present
	found := false
	for _, pattern := range updatedPatterns {
		if pattern == customPattern {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Custom pattern '%s' not found after adding", customPattern)
	}

	// Reset to defaults for other tests
	ResetToDefaults()
}

func TestRemoveExclusion(t *testing.T) {
	// Get initial patterns
	initialPatterns := GetDefaultExcludePatterns()
	initialCount := len(initialPatterns)

	// Remove a known pattern
	patternToRemove := ".git"
	RemoveExclusion(patternToRemove)

	// Check that it was removed
	updatedPatterns := GetDefaultExcludePatterns()
	if len(updatedPatterns) != initialCount-1 {
		t.Errorf("Expected %d patterns after removing, got %d", initialCount-1, len(updatedPatterns))
	}

	// Check that the pattern is not present
	for _, pattern := range updatedPatterns {
		if pattern == patternToRemove {
			t.Errorf("Pattern '%s' still found after removal", patternToRemove)
		}
	}

	// Reset to defaults for other tests
	ResetToDefaults()
}

func TestResetToDefaults(t *testing.T) {
	// Get initial state
	initialPatterns := GetDefaultExcludePatterns()
	initialCount := len(initialPatterns)

	// Modify the patterns
	AddCustomExclusion("*.custom1")
	AddCustomExclusion("*.custom2")
	RemoveExclusion(".git")

	// Verify modification
	modifiedPatterns := GetDefaultExcludePatterns()
	if len(modifiedPatterns) == initialCount {
		t.Error("Patterns should have been modified before reset")
	}

	// Reset to defaults
	ResetToDefaults()

	// Verify reset
	resetPatterns := GetDefaultExcludePatterns()
	if len(resetPatterns) != initialCount {
		t.Errorf("Expected %d patterns after reset, got %d", initialCount, len(resetPatterns))
	}

	// Check that .git is back
	found := false
	for _, pattern := range resetPatterns {
		if pattern == ".git" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Pattern '.git' should be present after reset")
	}
}

func TestShouldIncludeFileWithDefaults(t *testing.T) {
	defaultExclusions := GetDefaultExcludePatterns()

	testCases := []struct {
		path     string
		expected bool
		desc     string
	}{
		{"main.go", true, "Go source file should be included"},
		{"README.md", true, "README file should be included"},
		{"node_modules/package/index.js", false, "Files in node_modules should be excluded"},
		{".venv/lib/python3.9/site-packages/module.py", false, "Files in .venv should be excluded"},
		{"build/output.exe", false, "Files in build directory should be excluded"},
		{"src/main.exe", false, "Executable files should be excluded"},
		{"image.jpg", false, "Image files should be excluded"},
		{"video.mp4", false, "Video files should be excluded"},
		{"package-lock.json", false, "Lock files should be excluded"},
		{"temp.log", false, "Log files should be excluded"},
		{".DS_Store", false, "OS files should be excluded"},
		{"src/component.tsx", true, "TypeScript React files should be included"},
		{"docs/guide.md", true, "Documentation should be included"},
	}

	for _, tc := range testCases {
		result := ShouldIncludeFile(tc.path, nil, defaultExclusions)
		if result != tc.expected {
			t.Errorf("%s: expected %v, got %v for path '%s'", tc.desc, tc.expected, result, tc.path)
		}
	}
}
