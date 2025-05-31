# Download script for gingest binaries
# This script helps download specific platform binaries from GitHub releases

param(
    [Parameter(Mandatory=$false)]
    [ValidateSet("windows-amd64", "windows-arm64", "linux-amd64", "linux-arm64", "darwin-amd64", "darwin-arm64", "all")]
    [string]$Platform = "windows-amd64",
    
    [Parameter(Mandatory=$false)]
    [string]$Version = "latest",
    
    [Parameter(Mandatory=$false)]
    [string]$OutputDir = "."
)

$repo = "prashanth1k/gingest"
$baseUrl = "https://github.com/$repo"

Write-Host "Gingest Binary Downloader" -ForegroundColor Green
Write-Host "=========================" -ForegroundColor Green
Write-Host "Platform: $Platform" -ForegroundColor Yellow
Write-Host "Version: $Version" -ForegroundColor Yellow
Write-Host "Output Directory: $OutputDir" -ForegroundColor Yellow
Write-Host ""

# Create output directory if it doesn't exist
if (!(Test-Path $OutputDir)) {
    New-Item -ItemType Directory -Path $OutputDir -Force | Out-Null
    Write-Host "Created directory: $OutputDir" -ForegroundColor Blue
}

# Function to download a specific binary
function Download-Binary {
    param($platform, $version, $outputDir)
    
    $extension = if ($platform.StartsWith("windows")) { ".exe" } else { "" }
    $filename = "gingest-$platform$extension"
    $downloadUrl = if ($version -eq "latest") {
        "$baseUrl/releases/latest/download/$filename"
    } else {
        "$baseUrl/releases/download/$version/$filename"
    }
    
    $outputPath = Join-Path $outputDir $filename
    
    Write-Host "Downloading $filename..." -ForegroundColor Blue
    try {
        Invoke-WebRequest -Uri $downloadUrl -OutFile $outputPath -ErrorAction Stop
        Write-Host "✓ Downloaded: $outputPath" -ForegroundColor Green
        
        # Make executable on Unix-like systems (if running on WSL/Linux)
        if (!$platform.StartsWith("windows") -and (Get-Command chmod -ErrorAction SilentlyContinue)) {
            chmod +x $outputPath
            Write-Host "✓ Made executable: $outputPath" -ForegroundColor Green
        }
        
        return $true
    } catch {
        Write-Host "✗ Failed to download $filename`: $($_.Exception.Message)" -ForegroundColor Red
        return $false
    }
}

# Function to download checksums
function Download-Checksums {
    param($version, $outputDir)
    
    $downloadUrl = if ($version -eq "latest") {
        "$baseUrl/releases/latest/download/checksums.txt"
    } else {
        "$baseUrl/releases/download/$version/checksums.txt"
    }
    
    $outputPath = Join-Path $outputDir "checksums.txt"
    
    Write-Host "Downloading checksums..." -ForegroundColor Blue
    try {
        Invoke-WebRequest -Uri $downloadUrl -OutFile $outputPath -ErrorAction Stop
        Write-Host "✓ Downloaded: $outputPath" -ForegroundColor Green
        return $true
    } catch {
        Write-Host "✗ Failed to download checksums: $($_.Exception.Message)" -ForegroundColor Red
        return $false
    }
}

# Download based on platform selection
if ($Platform -eq "all") {
    Write-Host "Downloading all platforms..." -ForegroundColor Cyan
    $platforms = @("windows-amd64", "windows-arm64", "linux-amd64", "linux-arm64", "darwin-amd64", "darwin-arm64")
    $successCount = 0
    
    foreach ($p in $platforms) {
        if (Download-Binary $p $Version $OutputDir) {
            $successCount++
        }
    }
    
    # Download checksums
    Download-Checksums $Version $OutputDir | Out-Null
    
    Write-Host ""
    Write-Host "Downloaded $successCount of $($platforms.Count) binaries" -ForegroundColor $(if ($successCount -eq $platforms.Count) { "Green" } else { "Yellow" })
} else {
    # Download specific platform
    if (Download-Binary $Platform $Version $OutputDir) {
        Download-Checksums $Version $OutputDir | Out-Null
        Write-Host ""
        Write-Host "Download complete!" -ForegroundColor Green
    } else {
        Write-Host ""
        Write-Host "Download failed!" -ForegroundColor Red
        exit 1
    }
}

Write-Host ""
Write-Host "Usage examples:" -ForegroundColor Cyan
Write-Host "  .\download.ps1                                    # Download Windows AMD64 (default)"
Write-Host "  .\download.ps1 -Platform linux-amd64              # Download Linux AMD64"
Write-Host "  .\download.ps1 -Platform darwin-arm64             # Download macOS ARM64 (Apple Silicon)"
Write-Host "  .\download.ps1 -Platform all                      # Download all platforms"
Write-Host "  .\download.ps1 -Version v0.1.0                    # Download specific version"
Write-Host "  .\download.ps1 -Platform windows-amd64 -OutputDir ./bin  # Custom output directory" 