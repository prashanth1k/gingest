# Default Exclusions Demo

This document demonstrates the comprehensive default exclusions in gingest and how to customize them.

## Default Exclusions

By default, gingest excludes the following types of files and directories:

### Version Control

- `.git`, `.svn`, `.hg`, `.bzr`

### Dependencies by Language

**JavaScript/Node.js:**

- `node_modules`, `.npm`, `.yarn`, `bower_components`

**Python:**

- `.venv`, `venv`, `env`, `__pycache__`, `.pytest_cache`, `.mypy_cache`, `.tox`, `site-packages`, `dist`, `build`, `*.egg-info`

**Go:**

- `vendor`

**Java/JVM:**

- `target`, `.gradle`, `.m2`, `build`, `out`

**C#/.NET:**

- `bin`, `obj`, `packages`, `.nuget`

**Ruby:**

- `.bundle`, `vendor/bundle`, `.gem`

**PHP:**

- `vendor`, `composer.phar`

**Rust:**

- `target`, `Cargo.lock`

**C/C++:**

- `build`, `cmake-build-*`, `.cmake`

### IDE and Editor Files

- `.vscode`, `.idea`, `.vs`, `.atom`, `.sublime-*`

### OS Generated Files

- `.DS_Store`, `Thumbs.db`, `Desktop.ini`

### Binary and Media Files

- Executables: `*.exe`, `*.dll`, `*.so`, `*.dylib`
- Images: `*.jpg`, `*.png`, `*.gif`, `*.mp4`, `*.mp3`
- Archives: `*.zip`, `*.tar`, `*.rar`

### Lock Files

- `package-lock.json`, `yarn.lock`, `Pipfile.lock`, `poetry.lock`, `go.sum`

### Temporary Files

- `*.log`, `*.tmp`, `*.temp`, `*~`, `*.swp`

## Usage Examples

### Use Default Exclusions

```bash
# Uses comprehensive default exclusions
gingest --source=./my-project
```

### Add Custom Exclusions

```bash
# Add custom patterns to defaults
gingest --source=./my-project --exclude="*.custom,temp-*,debug/"
```

### Disable Default Exclusions

```bash
# Process all files (no exclusions)
gingest --source=./my-project --exclude=""
```

### Include Specific Files Despite Exclusions

```bash
# Include only Go and Markdown files, excluding everything else
gingest --source=./my-project --include="*.go,*.md"
```

### Mixed Include/Exclude

```bash
# Include Go files but exclude test files
gingest --source=./my-project --include="*.go" --exclude="*_test.go"
```

## Programmatic Usage

```go
package main

import (
    "fmt"
    "github.com/prashanth1k/gingest/internal/utils"
)

func main() {
    // Get all default exclusions
    defaults := utils.GetDefaultExcludePatterns()
    fmt.Printf("Total default exclusions: %d\n", len(defaults))

    // Get only directory exclusions
    dirs := utils.GetDefaultDirectoryExclusions()
    fmt.Printf("Directory exclusions: %d\n", len(dirs))

    // Get only file pattern exclusions
    files := utils.GetDefaultFileExclusions()
    fmt.Printf("File pattern exclusions: %d\n", len(files))

    // Add custom exclusion at runtime
    utils.AddCustomExclusion("*.myext")

    // Remove a default exclusion
    utils.RemoveExclusion(".git")

    // Reset to original defaults
    utils.ResetToDefaults()
}
```

## Benefits

1. **Zero Configuration**: Works out of the box for most projects
2. **Language Agnostic**: Covers popular languages and frameworks
3. **Customizable**: Easy to override or extend
4. **Performance**: Avoids processing unnecessary files
5. **Security**: Excludes environment files that may contain secrets
