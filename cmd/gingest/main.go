package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/prashanth1k/gingest/internal/ingester"
	"github.com/prashanth1k/gingest/internal/types"
	"github.com/prashanth1k/gingest/internal/utils"
)

// Version information - set during build with ldflags
var (
	Version   = "dev"     // Version number
	GitCommit = "unknown" // Git commit hash
	BuildDate = "unknown" // Build date
)

func printUsage() {
	fmt.Fprintf(os.Stderr, `gingest - Convert codebases into LLM-friendly text digests

USAGE:
    gingest --source=<path|url> [OPTIONS]

EXAMPLES:
    gingest --source=./my-project
    gingest --source=https://github.com/user/repo.git --output=repo.md
    gingest --source=./project --maxsize=1048576 --output=digest.md
    gingest --source=https://github.com/user/repo.git --branch=develop
    gingest --source=./project --exclude="*.log,node_modules,*.tmp"
    gingest --source=./project --include="*.go,*.md" --exclude=".git"

OPTIONS:
    --source=<path|url>    Source path (local directory or Git URL) [REQUIRED]
    --output=<file>        Output file path (default: digest.md)
    --branch=<name>        Target branch for Git repositories
    --maxsize=<bytes>      Maximum file size in bytes (default: 2MB)
    --exclude=<patterns>   Comma-separated exclude patterns (default: comprehensive list)
    --include=<patterns>   Comma-separated include patterns (overrides excludes)
    --version              Show version information
    --help, -h             Show this help message

DESCRIPTION:
    gingest processes local directories or remote Git repositories and generates
    consolidated, LLM-friendly text digests. The output contains all file contents
    with clear separators, optimized for Large Language Model consumption.

    Supports GitHub, GitLab, and other Git hosting services. Files exceeding the
    size limit are skipped with a descriptive message.

    Default exclusions include: dependency directories (.venv, venv, node_modules,
    vendor, target, build, etc.), version control (.git, .svn), IDE files (.vscode,
    .idea), OS files (.DS_Store, Thumbs.db), temporary files (*.tmp, *.log),
    binary files (*.exe, *.dll, *.so), media files (*.jpg, *.mp4, *.mp3),
    and many more. Use --exclude="" to disable defaults.

For more information, visit: https://github.com/prashanth1k/gingest
`)
}

func main() {
	// Define CLI flags
	var sourcePath = flag.String("source", "", "Source path (local directory or Git URL)")
	var outputFile = flag.String("output", "digest.md", "Output file path")
	var targetBranch = flag.String("branch", "", "Target branch for Git repositories")
	var maxFileSize = flag.Int64("maxsize", 2*1024*1024, "Maximum file size in bytes (default: 2MB)")
	var excludePatterns = flag.String("exclude", "", "Comma-separated exclude patterns")
	var includePatterns = flag.String("include", "", "Comma-separated include patterns")
	var showVersion = flag.Bool("version", false, "Show version information")

	// Set custom usage function
	flag.Usage = printUsage

	// Parse flags
	flag.Parse()

	// Handle version flag
	if *showVersion {
		fmt.Printf("gingest version %s\n", Version)
		fmt.Printf("Git commit: %s\n", GitCommit)
		fmt.Printf("Build date: %s\n", BuildDate)
		fmt.Printf("Go version: %s\n", runtime.Version())
		return
	}

	// Check if source is provided
	if *sourcePath == "" {
		fmt.Fprintf(os.Stderr, "Error: --source is required\n\n")
		flag.Usage()
		os.Exit(1)
	}

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Print parsed values
	fmt.Printf("Source Path: %s\n", *sourcePath)
	fmt.Printf("Output File: %s\n", *outputFile)
	fmt.Printf("Max File Size: %d bytes (%.1f MB)\n", *maxFileSize, float64(*maxFileSize)/(1024*1024))
	if *targetBranch != "" {
		fmt.Printf("Target Branch: %s\n", *targetBranch)
	}
	if *excludePatterns != "" {
		fmt.Printf("Exclude Patterns: %s\n", *excludePatterns)
	}
	if *includePatterns != "" {
		fmt.Printf("Include Patterns: %s\n", *includePatterns)
	}

	// Parse patterns
	var excludeList []string
	// Check if exclude flag was explicitly set
	excludeFlag := flag.Lookup("exclude")
	if excludeFlag != nil && excludeFlag.Value.String() != excludeFlag.DefValue {
		// Flag was explicitly set (even if to empty string)
		excludeList = utils.ParsePatterns(*excludePatterns)
	} else {
		// Flag was not set, use comprehensive default exclusions
		excludeList = utils.GetDefaultExcludePatterns()
	}
	includeList := utils.ParsePatterns(*includePatterns)

	var filesData []types.FileInfo
	var stats types.Stats
	var err error

	// Check if source is a Git URL
	if utils.IsGitURL(*sourcePath) {
		// Check if git is available
		if !utils.IsGitAvailable() {
			log.Fatal("Error: git command not found. Please install Git to process remote repositories.")
		}

		fmt.Printf("Processing remote Git repository: %s\n", *sourcePath)
		if *targetBranch != "" {
			fmt.Printf("Cloning branch: %s\n", *targetBranch)
		} else {
			fmt.Println("Cloning default branch...")
		}

		_, filesData, stats, err = ingester.ProcessRemoteRepoWithPatterns(*sourcePath, *targetBranch, *maxFileSize, includeList, excludeList)
		if err != nil {
			log.Fatalf("Error processing remote repository: %v", err)
		}
		fmt.Println("Clone successful.")
	} else if info, err := os.Stat(*sourcePath); err == nil && info.IsDir() {
		fmt.Printf("Processing local directory: %s\n", *sourcePath)
		fmt.Println("Scanning files...")

		filesData, stats, err = ingester.ProcessLocalDirectoryWithPatterns(*sourcePath, *maxFileSize, includeList, excludeList)
		if err != nil {
			log.Fatalf("Error processing directory: %v", err)
		}
	} else {
		log.Fatal("Source must be a valid local directory or Git URL")
	}

	fmt.Printf("Found %d files:\n", len(filesData))
	for _, fileInfo := range filesData {
		if fileInfo.Error != nil {
			fmt.Printf("  %s (ERROR: %v)\n", fileInfo.RelativePath, fileInfo.Error)
		} else {
			fmt.Printf("  %s (%d bytes)\n", fileInfo.RelativePath, len(fileInfo.Content))
		}
	}

	// Write digest to output file
	fmt.Printf("Writing digest to %s...\n", *outputFile)
	err = ingester.WriteDigest(*outputFile, filesData, stats)
	if err != nil {
		log.Fatalf("Error writing digest: %v", err)
	}

	fmt.Printf("Digest created: %s\n", *outputFile)

	// Print summary to stdout
	fmt.Println("\n" + utils.GenerateSummaryString(stats))
}
