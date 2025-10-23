#!/bin/bash

# Got - Go Reverse Tunnel Installer
# Usage: curl -sSL https://raw.githubusercontent.com/HeyRistaa/got/main/install.sh | bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Detect OS and architecture
OS=""
ARCH=""

case "$(uname -s)" in
    Linux*)     OS="linux";;
    Darwin*)    OS="darwin";;
    CYGWIN*)    OS="windows";;
    MINGW*)     OS="windows";;
    *)          echo -e "${RED}Unsupported OS: $(uname -s)${NC}"; exit 1;;
esac

case "$(uname -m)" in
    x86_64)     ARCH="amd64";;
    arm64)      ARCH="arm64";;
    aarch64)    ARCH="arm64";;
    *)          echo -e "${RED}Unsupported architecture: $(uname -m)${NC}"; exit 1;;
esac

# Set version (update this for new releases)
VERSION="v0.0.2"

# Determine binary name
if [ "$OS" = "windows" ]; then
    BINARY_NAME="got-windows-${ARCH}.exe"
else
    BINARY_NAME="got-${OS}-${ARCH}"
fi

# Download URL
DOWNLOAD_URL="https://github.com/HeyRistaa/got/releases/download/${VERSION}/${BINARY_NAME}"

echo -e "${GREEN}Installing got ${VERSION}...${NC}"
echo -e "${YELLOW}OS: ${OS}${NC}"
echo -e "${YELLOW}Architecture: ${ARCH}${NC}"
echo -e "${YELLOW}Binary: ${BINARY_NAME}${NC}"

# Create temp directory
TEMP_DIR=$(mktemp -d)
cd "$TEMP_DIR"

# Download binary
echo -e "${YELLOW}Downloading...${NC}"
if command -v curl >/dev/null 2>&1; then
    curl -L -o "$BINARY_NAME" "$DOWNLOAD_URL"
elif command -v wget >/dev/null 2>&1; then
    wget -O "$BINARY_NAME" "$DOWNLOAD_URL"
else
    echo -e "${RED}Neither curl nor wget found. Please install one of them.${NC}"
    exit 1
fi

# Make executable
chmod +x "$BINARY_NAME"

# Install to /usr/local/bin (requires sudo)
echo -e "${YELLOW}Installing to /usr/local/bin...${NC}"
if sudo mv "$BINARY_NAME" /usr/local/bin/got; then
    echo -e "${GREEN}✅ Successfully installed got to /usr/local/bin/got${NC}"
    echo -e "${GREEN}Run 'got --help' to get started!${NC}"
else
    echo -e "${YELLOW}⚠️  Could not install to /usr/local/bin (permission denied)${NC}"
    echo -e "${YELLOW}You can manually move the binary:${NC}"
    echo -e "${YELLOW}  sudo mv ${TEMP_DIR}/${BINARY_NAME} /usr/local/bin/got${NC}"
    echo -e "${YELLOW}  chmod +x /usr/local/bin/got${NC}"
fi

# Cleanup
rm -rf "$TEMP_DIR"
