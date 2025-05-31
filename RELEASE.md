# Release Guide for gingest

## Go Versioning Conventions

### 1. **Git Tags (Primary Method)**

Go uses **Git tags** for versioning following [Semantic Versioning](https://semver.org/):

```bash
# Major version (breaking changes)
git tag v2.0.0

# Minor version (new features, backward compatible)
git tag v1.1.0

# Patch version (bug fixes)
git tag v1.0.1

# Pre-release versions
git tag v1.0.0-alpha.1
git tag v1.0.0-beta.1
git tag v1.0.0-rc.1
```

### 2. **Version Information in Code**

The CLI includes version information that gets set during build:

```go
// In cmd/gingest/main.go
var (
    Version   = "dev"      // Set via -ldflags during build
    GitCommit = "unknown"  // Git commit hash
    BuildDate = "unknown"  // Build timestamp
)
```

### 3. **go.mod Version**

Specifies minimum Go version requirement:

```go
module github.com/prashanth1k/gingest

go 1.21
```

## Release Process

### Prerequisites

1. Initialize Git repository:

   ```bash
   git init
   git add .
   git commit -m "Initial commit"
   git branch -M main
   git remote add origin https://github.com/prashanth1k/gingest.git
   git push -u origin main
   ```

2. Set up GitHub repository secrets (for Docker publishing):
   - `DOCKER_USERNAME`: Your Docker Hub username
   - `DOCKER_PASSWORD`: Your Docker Hub password/token

### Automated Release (Recommended)

1. **Create and push a version tag:**

   ```bash
   # Create a new version tag
   git tag v1.0.0

   # Push the tag to GitHub
   git push origin v1.0.0
   ```

2. **GitHub Actions automatically:**
   - Runs all tests (unit + integration)
   - Builds binaries for multiple platforms:
     - Linux (amd64, arm64)
     - macOS (amd64, arm64)
     - Windows (amd64, arm64)
   - Creates SHA256 checksums
   - Creates GitHub release with:
     - Release notes
     - Binary downloads
     - Installation instructions
   - Builds and pushes Docker image to Docker Hub

### Manual Release (Using Makefile)

1. **Build locally:**

   ```bash
   # Build for current platform
   make build VERSION=v1.0.0

   # Build for all platforms
   make build-all VERSION=v1.0.0

   # Test the build
   make run-version
   ```

2. **Create release manually:**
   ```bash
   make release VERSION=v1.0.0
   ```

### Development Builds

```bash
# Build with development version
make build

# Build with custom version
make build VERSION=v1.0.0-dev

# Build Docker image
make docker-build VERSION=v1.0.0
```

## Version Examples

### Semantic Versioning

- `v1.0.0` - First stable release
- `v1.1.0` - New features added
- `v1.1.1` - Bug fixes
- `v2.0.0` - Breaking changes

### Pre-release Versions

- `v1.0.0-alpha.1` - Alpha release
- `v1.0.0-beta.1` - Beta release
- `v1.0.0-rc.1` - Release candidate

## Installation Methods

### For Users

1. **Go install (latest):**

   ```bash
   go install github.com/prashanth1k/gingest/cmd/gingest@latest
   ```

2. **Go install (specific version):**

   ```bash
   go install github.com/prashanth1k/gingest/cmd/gingest@v1.0.0
   ```

3. **Download binary:**

   - Visit GitHub releases page
   - Download appropriate binary for your platform
   - Verify with checksums if needed

4. **Docker:**
   ```bash
   docker pull prashanth1k/gingest:latest
   docker pull prashanth1k/gingest:v1.0.0
   ```

### For Developers

```bash
# Clone and build
git clone https://github.com/prashanth1k/gingest.git
cd gingest
make build

# Or install from source
make install
```

## GitHub Actions Workflows

### 1. **CI Workflow** (`.github/workflows/ci.yml`)

- Triggers: Push to `main`/`develop`, Pull Requests
- Runs on: Ubuntu, Windows, macOS
- Go versions: 1.21, 1.22
- Actions:
  - Unit tests
  - Integration tests
  - Race condition tests
  - Code formatting checks
  - Static analysis (staticcheck)
  - Cross-platform builds
  - Coverage reporting

### 2. **Release Workflow** (`.github/workflows/release.yml`)

- Triggers: Git tags matching `v*`
- Actions:
  - Run full test suite
  - Build binaries for all platforms
  - Create GitHub release
  - Build and push Docker images
  - Generate checksums

## Build Flags

The build process uses Go's `-ldflags` to inject version information:

```bash
go build -ldflags "-X main.Version=v1.0.0 -X main.GitCommit=abc123 -X main.BuildDate=2024-01-01T12:00:00Z"
```

## Docker Images

Images are published to Docker Hub:

- `prashanth1k/gingest:latest` - Latest release
- `prashanth1k/gingest:v1.0.0` - Specific version

Usage:

```bash
# Run with local directory
docker run --rm -v $(pwd):/workspace -v $(pwd)/output:/output prashanth1k/gingest:latest --source=/workspace --output=/output/digest.md

# Run with Git repository
docker run --rm -v $(pwd)/output:/output prashanth1k/gingest:latest --source=https://github.com/user/repo.git --output=/output/digest.md
```

## Troubleshooting

### Common Issues

1. **Git not initialized:**

   ```bash
   fatal: not a git repository
   ```

   Solution: Initialize git repository first

2. **Missing Docker secrets:**

   ```
   Error: Username and password required
   ```

   Solution: Add `DOCKER_USERNAME` and `DOCKER_PASSWORD` to GitHub secrets

3. **Build failures:**

   - Check Go version compatibility
   - Ensure all tests pass locally
   - Verify import paths are correct

4. **"Dependencies file is not found" warning:**
   ```
   Warning: Restore cache failed: Dependencies file is not found in /home/runner/work/gingest/gingest. Supported file pattern: go.sum
   ```
   Solution: This warning is harmless and has been fixed in the workflows. It occurs because gingest has no external dependencies (uses only Go standard library), so no `go.sum` file is created. The workflows now disable automatic caching to prevent this warning.

### Testing Releases

Before creating a release:

```bash
# Run all tests
make test-all

# Check formatting and linting
make fmt lint staticcheck

# Test build for all platforms
make build-all

# Test Docker build
make docker-build
```

## Best Practices

1. **Always test before releasing**
2. **Use semantic versioning**
3. **Write meaningful release notes**
4. **Test installation methods**
5. **Keep CHANGELOG.md updated**
6. **Tag releases consistently**
7. **Verify checksums for security**

## Next Steps

Once you initialize Git:

1. Push your code to GitHub
2. Set up repository secrets for Docker
3. Create your first release: `git tag v1.0.0 && git push origin v1.0.0`
4. Watch the GitHub Actions build and release automatically!

```

```
