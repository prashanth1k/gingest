# gingest

[![Go Reference](https://pkg.go.dev/badge/github.com/prashanth1k/gingest.svg)](https://pkg.go.dev/github.com/prashanth1k/gingest)
[![Go Report Card](https://goreportcard.com/badge/github.com/prashanth1k/gingest)](https://goreportcard.com/report/github.com/prashanth1k/gingest)

**gingest** is a command-line interface (CLI) tool and Go library designed to process local file directories or remote Git repositories (GitHub, GitLab) and generate consolidated, LLM-friendly text digests of their codebase. This digest is optimized for providing context to Large Language Models, enabling tasks like code analysis, summarization, and question-answering.

## Features

- **Local Directory Processing**: Recursively processes local directories
- **Remote Git Repository Support**: Clones and processes GitHub/GitLab repositories
- **Branch Selection**: Specify target branch for Git repositories
- **File Size Filtering**: Skip files exceeding configurable size limits (default: 2MB)
- **Binary File Detection**: Automatically detects and handles binary files
- **Jupyter Notebook Support**: Extracts content from `.ipynb` files
- **Include/Exclude Patterns**: Filter files using glob patterns
- **README Prioritization**: README files appear first in the digest
- **Directory Tree Output**: Visual directory structure in the digest
- **Summary Statistics**: Detailed processing statistics and metadata
- **Concurrent Processing**: Fast file processing using goroutines
- **Structured Output**: Generates LLM-friendly text digests with clear file separators
- **Error Handling**: Graceful handling of file read errors and Git clone failures
- **Go Library**: Use as a library in your Go applications
- **Cross-platform**: Works on Windows, macOS, and Linux

## Installation

### As a CLI tool

```bash
go install github.com/prashanth1k/gingest/cmd/gingest@latest
```

### As a Go library

```bash
go get github.com/prashanth1k/gingest
```

## Usage

### CLI Usage

#### Process a local directory

```bash
gingest --source=./my-project --output=digest.md
```

#### Process a remote Git repository

```bash
gingest --source=https://github.com/user/repo.git --output=repo_digest.md
```

#### Process with include/exclude patterns

```bash
# Include only Go and Python files
gingest --source=./project --include="*.go,*.py" --output=code_digest.md

# Exclude test files and documentation
gingest --source=./project --exclude="*_test.go,*.md,docs/*" --output=digest.md
```

#### Process specific branch with size limit

```bash
gingest --source=https://github.com/user/repo.git --branch=develop --maxsize=1048576 --output=digest.md
```

#### CLI Flags

- `--source`: Source path (local directory or Git URL) **[required]**
- `--output`: Output file path (default: `digest.md`)
- `--branch`: Target branch for Git repositories (optional)
- `--maxsize`: Maximum file size in bytes (default: 2MB)
- `--include`: Comma-separated glob patterns for files to include (optional)
- `--exclude`: Comma-separated glob patterns for files to exclude (optional)

**Default exclude patterns**: `.git`, `*.log`, `node_modules`, `.DS_Store`, `Thumbs.db`, `*.tmp`, `*.temp`, `.vscode`, `.idea`, `*.swp`, `*.swo`, `*~`

### Library Usage

```go
package main

import (
    "fmt"
    "log"

    "github.com/prashanth1k/gingest"
)

func main() {
    // Process local directory
    config := gingest.Config{
        Source:      "./my-project",
        OutputFile:  "digest.md",
        MaxFileSize: 1024 * 1024, // 1MB limit
    }

    err := gingest.ProcessAndWriteDigest(config)
    if err != nil {
        log.Fatalf("Error: %v", err)
    }

    fmt.Println("Digest created successfully!")
}
```

#### Advanced Usage with Statistics

```go
// Get file information and statistics
filesData, stats, err := gingest.ProcessCodebaseWithStats(config)
if err != nil {
    log.Fatalf("Error: %v", err)
}

fmt.Printf("Processed %d files, %d directories\n", stats.NumFilesProcessed, stats.NumDirsProcessed)
fmt.Printf("Binary files: %d, Skipped files: %d\n", stats.NumBinaryFiles, stats.NumSkippedFiles)
fmt.Printf("Total content: %.2f KB\n", float64(stats.TotalContentBytes)/1024)

// Process each file
for _, fileInfo := range filesData {
    if fileInfo.Error != nil {
        fmt.Printf("Error reading %s: %v\n", fileInfo.RelativePath, fileInfo.Error)
        continue
    }

    if fileInfo.IsBinary {
        fmt.Printf("Binary file: %s\n", fileInfo.RelativePath)
    } else {
        fmt.Printf("Text file: %s (%d bytes)\n", fileInfo.RelativePath, len(fileInfo.Content))
    }
}
```

## Output Format

The generated digest includes:

1. **Summary Section**: Processing statistics and metadata
2. **Directory Tree**: Visual representation of the project structure
3. **File Contents**: Individual file contents with clear separators

```markdown
# Codebase Digest Summary

**Source:** ./my-project
**Generated:** 2024-01-15 10:30:45

## Statistics

- **Total Files:** 25
- **Directories:** 8
- **Binary Files:** 3
- **Skipped Files:** 1
- **Total Content Size:** 45.67 KB

---

## Directory Structure
```

my-project/
├── README.md
├── main.go
├── config/
│ └── config.go
└── utils/
├── helper.go
└── binary.bin (Binary)

```

---

================================================
FILE: README.md
================================================
# My Project

This is a sample project...

================================================
FILE: main.go
================================================
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}

```

**Features of the output:**

- README files appear first
- Files sorted alphabetically within their categories
- Binary files marked as `[Binary File]`
- Large files marked as `[File content skipped: Exceeds max size]`
- Jupyter notebooks parsed and formatted with cell structure

## Configuration

### Config Struct

```go
type Config struct {
    Source         string   // Local directory path or Git repository URL
    OutputFile     string   // Output file path for the digest
    TargetBranch   string   // Target branch for Git repositories (optional)
    MaxFileSize    int64    // Maximum file size in bytes (0 = no limit)
    IncludePatterns []string // Glob patterns for files to include (optional)
    ExcludePatterns []string // Glob patterns for files to exclude (optional)
}
```

## Examples

See the [examples](./examples/) directory for complete usage examples:

- [Basic Usage](./examples/basic/main.go) - Simple local and remote processing

## Development

### Building from Source

```bash
git clone https://github.com/prashanth1k/gingest.git
cd gingest
go build -o gingest cmd/gingest/main.go
```

### Running Tests

```bash
# Unit tests
go test ./...

# Integration tests
go test -tags=integration -v

# All tests
go test -tags=integration ./...
```

### Project Structure

```
gingest/
├── cmd/gingest/              # CLI application
├── internal/                 # Internal packages
│   ├── ingester/            # Core processing logic
│   ├── notebookparser/      # Jupyter notebook parsing
│   ├── types/               # Type definitions
│   └── utils/               # Utility functions
├── examples/                # Usage examples
├── gingest.go               # Public API
├── gingest_test.go          # Unit tests
├── integration_test.go      # Integration tests
└── README.md
```

## Supported File Types

- **Text Files**: `.go`, `.py`, `.js`, `.ts`, `.java`, `.cpp`, `.c`, `.h`, `.md`, `.txt`, `.yaml`, `.yml`, `.json`, `.xml`, `.html`, `.css`, `.sql`, etc.
- **Jupyter Notebooks**: `.ipynb` files are parsed to extract markdown and code cells
- **Binary Files**: Detected automatically and marked as `[Binary File]`
- **Large Files**: Files exceeding size limit are marked as `[File content skipped]`

## Performance

- **Concurrent Processing**: Files are processed concurrently using goroutines for improved performance
- **Memory Efficient**: Streams file content without loading entire codebase into memory
- **Fast Git Operations**: Uses shallow clones (`--depth 1`) for remote repositories

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Goals

- **Simple and Effective**: Provide an easy way to convert codebases into LLM-friendly format
- **Flexible Filtering**: Control which files are included based on size, patterns, and branch
- **User-Friendly**: Clear feedback and intuitive commands
- **Performance**: Fast and efficient processing leveraging Go's strengths
- **Cross-Platform**: Works consistently across different operating systems
- **Smart Processing**: Intelligently handle common repository elements like README files
