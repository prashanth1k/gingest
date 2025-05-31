package utils

// DefaultExclusions contains default exclusion patterns for various environments
// This follows Go's idiomatic approach of using package-level variables for configuration
// that can be modified at runtime if needed.
var DefaultExclusions = struct {
	// Directories to exclude by default
	Directories []string
	// File patterns to exclude by default
	FilePatterns []string
	// Combined patterns for convenience
	All []string
}{
	Directories: []string{
		// Version control
		".git",
		".svn",
		".hg",
		".bzr",

		// Dependencies - JavaScript/Node.js
		"node_modules",
		".npm",
		".yarn",
		"bower_components",

		// Dependencies - Python
		".venv",
		"venv",
		"env",
		".env",
		"__pycache__",
		".pytest_cache",
		".mypy_cache",
		".tox",
		"site-packages",
		"dist",
		"build",
		"*.egg-info",
		".coverage",

		// Dependencies - Go
		"vendor",

		// Dependencies - Java/JVM
		"target",
		".gradle",
		".m2",
		"build",
		"out",

		// Dependencies - .NET
		"bin",
		"obj",
		"packages",
		".nuget",

		// Dependencies - Ruby
		".bundle",
		"vendor/bundle",
		".gem",

		// Dependencies - PHP
		"vendor",
		"composer.phar",

		// Dependencies - Rust
		"target",

		// Dependencies - C/C++
		"build",
		"cmake-build-*",
		".cmake",

		// IDE and editor directories
		".vscode",
		".idea",
		".vs",
		".atom",
		".sublime-*",

		// OS generated
		".DS_Store",
		"Thumbs.db",
		"Desktop.ini",

		// Temporary and cache directories
		"tmp",
		"temp",
		".tmp",
		".temp",
		"cache",
		".cache",

		// Logs
		"logs",
		"log",

		// Documentation build outputs
		"_site",
		"docs/_build",
		".docusaurus",
		".next",
		".nuxt",
		"dist",
		"public",

		// Testing
		"coverage",
		".nyc_output",
		"test-results",
		"test-reports",

		// Docker
		".docker",

		// Terraform
		".terraform",
		"*.tfstate*",

		// Kubernetes
		".kube",
	},

	FilePatterns: []string{
		// Logs
		"*.log",
		"*.log.*",

		// Temporary files
		"*.tmp",
		"*.temp",
		"*~",
		"*.swp",
		"*.swo",
		".#*",
		"#*#",

		// OS files
		".DS_Store",
		"Thumbs.db",
		"Desktop.ini",
		"*.lnk",

		// Compiled files
		"*.o",
		"*.obj",
		"*.exe",
		"*.dll",
		"*.so",
		"*.dylib",
		"*.class",
		"*.pyc",
		"*.pyo",
		"*.pyd",

		// Archives
		"*.zip",
		"*.tar",
		"*.tar.gz",
		"*.tgz",
		"*.rar",
		"*.7z",

		// Images (usually not useful for code analysis)
		"*.jpg",
		"*.jpeg",
		"*.png",
		"*.gif",
		"*.bmp",
		"*.ico",
		"*.svg",
		"*.webp",

		// Videos
		"*.mp4",
		"*.avi",
		"*.mov",
		"*.wmv",
		"*.flv",
		"*.webm",

		// Audio
		"*.mp3",
		"*.wav",
		"*.flac",
		"*.aac",
		"*.ogg",

		// Documents (usually not code)
		"*.pdf",
		"*.doc",
		"*.docx",
		"*.xls",
		"*.xlsx",
		"*.ppt",
		"*.pptx",

		// Lock files (generated)
		"package-lock.json",
		"yarn.lock",
		"pnpm-lock.yaml",
		"bun.lockb",
		"Pipfile.lock",
		"poetry.lock",
		"pdm.lock",
		"Gemfile.lock",
		"composer.lock",
		"go.sum",
		"Cargo.lock",
		"mix.lock",
		"packages.lock.json",
		"project.assets.json",
		"*.lock",
		"flake.lock",
		"deno.lock",
		"shrinkwrap.yaml",
		"npm-shrinkwrap.json",
		"uv.lock",

		// Environment files (may contain secrets)
		".env",
		".env.*",
		"*.env",

		// IDE files
		"*.iml",
		"*.ipr",
		"*.iws",
		".project",
		".classpath",
		".settings",

		// Coverage reports
		"coverage.xml",
		"coverage.json",
		"lcov.info",
		".coverage",

		// Minified files
		"*.min.js",
		"*.min.css",

		// Source maps
		"*.map",
		"*.js.map",
		"*.css.map",
	},
}

// init function runs when the package is imported and combines all patterns
func init() {
	DefaultExclusions.All = make([]string, 0, len(DefaultExclusions.Directories)+len(DefaultExclusions.FilePatterns))
	DefaultExclusions.All = append(DefaultExclusions.All, DefaultExclusions.Directories...)
	DefaultExclusions.All = append(DefaultExclusions.All, DefaultExclusions.FilePatterns...)
}

// GetDefaultExcludePatterns returns all default exclusion patterns
// This maintains backward compatibility with existing code
func GetDefaultExcludePatterns() []string {
	return DefaultExclusions.All
}

// GetDefaultDirectoryExclusions returns only directory exclusion patterns
func GetDefaultDirectoryExclusions() []string {
	return DefaultExclusions.Directories
}

// GetDefaultFileExclusions returns only file pattern exclusions
func GetDefaultFileExclusions() []string {
	return DefaultExclusions.FilePatterns
}

// AddCustomExclusion allows runtime modification of default exclusions
// This is Go's idiomatic way to allow configuration changes
func AddCustomExclusion(pattern string) {
	DefaultExclusions.All = append(DefaultExclusions.All, pattern)
}

// RemoveExclusion removes a pattern from default exclusions
func RemoveExclusion(pattern string) {
	for i, p := range DefaultExclusions.All {
		if p == pattern {
			DefaultExclusions.All = append(DefaultExclusions.All[:i], DefaultExclusions.All[i+1:]...)
			break
		}
	}
}

// ResetToDefaults resets exclusions to the original defaults
func ResetToDefaults() {
	DefaultExclusions.All = make([]string, 0, len(DefaultExclusions.Directories)+len(DefaultExclusions.FilePatterns))
	DefaultExclusions.All = append(DefaultExclusions.All, DefaultExclusions.Directories...)
	DefaultExclusions.All = append(DefaultExclusions.All, DefaultExclusions.FilePatterns...)
}
