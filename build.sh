#!/usr/bin/env bash
set -e

VERSION="v0.0.3"
BUILD_DIR="releases"
DATE=$(date +%Y-%m-%d)

echo "🚀 Building Mein Tunnel $VERSION"
echo "📅 Build date: $DATE"
echo ""

# Clean previous builds
rm -rf $BUILD_DIR
mkdir -p $BUILD_DIR

echo "📦 Building binaries for all platforms..."
echo ""

# Build function
build_binary() {
    local os=$1
    local arch=$2
    local ext=$3
    
    echo "🔨 Building for $os/$arch..."
    
    GOOS=$os GOARCH=$arch go build -ldflags="-s -w" -o $BUILD_DIR/got-$os-$arch$ext ./cmd/client
    GOOS=$os GOARCH=$arch go build -ldflags="-s -w" -o $BUILD_DIR/server-$os-$arch$ext ./cmd/server
    
    echo "✅ Built got-$os-$arch$ext and server-$os-$arch$ext"
}

# Build for all platforms
build_binary "linux" "amd64" ""
build_binary "linux" "arm64" ""
build_binary "darwin" "amd64" ""
build_binary "darwin" "arm64" ""
build_binary "windows" "amd64" ".exe"
build_binary "windows" "arm64" ".exe"

echo ""
echo "🎉 Build complete!"
echo "📁 Binaries are in the $BUILD_DIR/ directory:"
ls -lh $BUILD_DIR/

echo ""
echo "✅ Ready for release v0.0.3!"

