#!/bin/bash
# Test installation script for v0.0.3

echo "🧪 Testing installation script..."
echo ""

# Test if binaries exist
echo "📋 Checking binaries..."
ls -lh releases/ || echo "❌ No releases directory found"

echo ""
echo "✅ To test the installation:"
echo "  1. Run: curl -sSL https://raw.githubusercontent.com/HeyRistaa/got/main/install.sh | bash"
echo "  2. The script will download the correct binary for your platform"
echo "  3. Check if 'got --help' works"

echo ""
echo "📝 Manual test steps:"
echo "  1. Make repository public"
echo "  2. Create GitHub release v0.0.3"
echo "  3. Upload all binaries from releases/ directory"
echo "  4. Test installation script from a clean machine"

