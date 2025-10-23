#!/bin/bash

# Build script for Mein Tunnel v0.0.2
set -e

VERSION="v0.0.2"
BUILD_DIR="releases"
DATE=$(date +%Y-%m-%d)

echo "🚀 Building Mein Tunnel $VERSION"
echo "📅 Build date: $DATE"
echo ""

# Clean previous builds
rm -rf $BUILD_DIR
mkdir -p $BUILD_DIR

# Build function
build_binary() {
    local os=$1
    local arch=$2
    local ext=$3
    
    echo "🔨 Building for $os/$arch..."
    
    # Set GOOS and GOARCH
    export GOOS=$os
    export GOARCH=$arch
    
    # Build client
    go build -ldflags="-s -w -X main.version=$VERSION" -o $BUILD_DIR/got-$os-$arch$ext ./cmd/client
    
    # Build server
    go build -ldflags="-s -w -X main.version=$VERSION" -o $BUILD_DIR/server-$os-$arch$ext ./cmd/server
    
    echo "✅ Built got-$os-$arch$ext and server-$os-$arch$ext"
}

# Build for all platforms
echo "📦 Building binaries for all platforms..."

# Linux
build_binary "linux" "amd64" ""
build_binary "linux" "arm64" ""

# macOS
build_binary "darwin" "amd64" ""
build_binary "darwin" "arm64" ""

# Windows
build_binary "windows" "amd64" ".exe"
build_binary "windows" "arm64" ".exe"

echo ""
echo "🎉 Build complete!"
echo "📁 Binaries are in the $BUILD_DIR/ directory:"
ls -la $BUILD_DIR/

echo ""
echo "📋 Release checklist:"
echo "  ✅ Binaries built for all platforms"
echo "  ✅ Version: $VERSION"
echo "  ✅ Build date: $DATE"
echo ""
echo "🚀 Ready for release!"
