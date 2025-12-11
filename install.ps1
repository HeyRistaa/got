# Got - Go Reverse Tunnel Installer for Windows
# Usage: Invoke-WebRequest -Uri "https://raw.githubusercontent.com/HeyRistaa/got/main/install.ps1" | Invoke-Expression

$ErrorActionPreference = "Stop"

# Set version (update this for new releases)
$VERSION = "v0.0.3"

# Detect architecture
$ARCH = if ([Environment]::Is64BitOperatingSystem) { "amd64" } else { "386" }
$BINARY_NAME = "got-windows-${ARCH}.exe"
$DOWNLOAD_URL = "https://github.com/HeyRistaa/got/releases/download/${VERSION}/${BINARY_NAME}"

Write-Host "Installing got ${VERSION}..." -ForegroundColor Green
Write-Host "Architecture: ${ARCH}" -ForegroundColor Yellow
Write-Host "Binary: ${BINARY_NAME}" -ForegroundColor Yellow

# Create temp directory
$TEMP_DIR = [System.IO.Path]::GetTempPath()
$TEMP_FILE = Join-Path $TEMP_DIR $BINARY_NAME

# Download binary
Write-Host "Downloading..." -ForegroundColor Yellow
Invoke-WebRequest -Uri $DOWNLOAD_URL -OutFile $TEMP_FILE

# Install to a directory in PATH
$INSTALL_DIR = "$env:USERPROFILE\bin"
if (-not (Test-Path $INSTALL_DIR)) {
    New-Item -ItemType Directory -Path $INSTALL_DIR -Force | Out-Null
}

$INSTALL_PATH = Join-Path $INSTALL_DIR "got.exe"
Move-Item $TEMP_FILE $INSTALL_PATH -Force

Write-Host "✅ Successfully installed got to ${INSTALL_PATH}" -ForegroundColor Green
Write-Host "Add ${INSTALL_DIR} to your PATH if it's not already there." -ForegroundColor Yellow
Write-Host "Run 'got --help' to get started!" -ForegroundColor Green
