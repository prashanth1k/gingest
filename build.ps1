# Build script for gingest
# This creates the same artifacts that GitHub Actions would create

Write-Host "Building gingest artifacts..." -ForegroundColor Green

# Get build information
$VERSION = "v0.1.0-local"
$COMMIT = git rev-parse --short HEAD
$DATE = Get-Date -Format "yyyy-MM-ddTHH:mm:ssZ"

Write-Host "Version: $VERSION" -ForegroundColor Yellow
Write-Host "Commit: $COMMIT" -ForegroundColor Yellow
Write-Host "Date: $DATE" -ForegroundColor Yellow

$LDFLAGS = "-X main.Version=$VERSION -X main.GitCommit=$COMMIT -X main.BuildDate=$DATE"

# Clean previous builds
Write-Host "Cleaning previous builds..." -ForegroundColor Blue
Remove-Item -Path "gingest*" -Force -ErrorAction SilentlyContinue

# Build for current platform first
Write-Host "Building for current platform..." -ForegroundColor Blue
go build -ldflags $LDFLAGS -o gingest.exe cmd/gingest/main.go

# Test the build
Write-Host "Testing build..." -ForegroundColor Blue
./gingest.exe --version

# Build for multiple platforms
Write-Host "Building for multiple platforms..." -ForegroundColor Blue

# Linux AMD64
$env:GOOS = "linux"
$env:GOARCH = "amd64"
go build -ldflags $LDFLAGS -o gingest-linux-amd64 cmd/gingest/main.go

# Linux ARM64
$env:GOARCH = "arm64"
go build -ldflags $LDFLAGS -o gingest-linux-arm64 cmd/gingest/main.go

# macOS AMD64
$env:GOOS = "darwin"
$env:GOARCH = "amd64"
go build -ldflags $LDFLAGS -o gingest-darwin-amd64 cmd/gingest/main.go

# macOS ARM64 (Apple Silicon)
$env:GOARCH = "arm64"
go build -ldflags $LDFLAGS -o gingest-darwin-arm64 cmd/gingest/main.go

# Windows AMD64
$env:GOOS = "windows"
$env:GOARCH = "amd64"
go build -ldflags $LDFLAGS -o gingest-windows-amd64.exe cmd/gingest/main.go

# Windows ARM64
$env:GOARCH = "arm64"
go build -ldflags $LDFLAGS -o gingest-windows-arm64.exe cmd/gingest/main.go

# Reset environment
Remove-Item Env:GOOS -ErrorAction SilentlyContinue
Remove-Item Env:GOARCH -ErrorAction SilentlyContinue

# List created files
Write-Host "Build artifacts created:" -ForegroundColor Green
Get-ChildItem -Name "gingest*" | ForEach-Object {
    $size = (Get-Item $_).Length
    $sizeKB = [math]::Round($size / 1KB, 2)
    Write-Host "  $_ ($sizeKB KB)" -ForegroundColor Cyan
}

# Generate checksums
Write-Host "Generating checksums..." -ForegroundColor Blue
$checksums = @()
Get-ChildItem "gingest*" | ForEach-Object {
    $hash = Get-FileHash $_.Name -Algorithm SHA256
    $checksums += "$($hash.Hash.ToLower())  $($_.Name)"
}
$checksums | Out-File -FilePath "checksums.txt" -Encoding UTF8
Write-Host "Checksums saved to checksums.txt" -ForegroundColor Green

Write-Host "Build complete! All artifacts are ready." -ForegroundColor Green
Write-Host "You can now distribute these binaries or create a release." -ForegroundColor Yellow 