# Changelog

All notable changes to this project will be documented in this file.

## [v0.0.3] - 2025-01-21

### 🧹 **Cleanup & Polish**
- **Removed ngrok references** - Cleaned up all ngrok comparisons for professional presentation
- **Updated terminology** - Changed "expose" to "tunnel" for better clarity
- **Improved branding** - Updated title to "Go Reverse Tunnel" for better recognition
- **Professional release** - Prepared for public repository release

### 🔒 **Security & Abuse Prevention**
- **Rate limiting** - Prevent abuse with configurable limits (5 tunnels per minute, 20 per hour)
- **IP-based tracking** - Track and limit connections by IP address
- **Removed hardcoded server** - No longer hardcodes production server IP in public repository
- **Explicit server specification** - Users must specify server via `-server` flag or `GOT_SERVER_HOST` env
- **Better error messages** - Clear guidance when server is not specified

### 🎨 **UI/UX Improvements**
- **Enhanced CLI colors** - Beautiful colored output with emojis and better visual feedback
- **Improved error messages** - More descriptive and user-friendly error handling
- **Better startup messages** - Clear visual indicators for tunnel establishment

### 🔧 **Technical Improvements**
- **Code organization** - Better structured codebase with clear separation of concerns
- **Documentation updates** - Comprehensive README with architecture diagrams
- **Version management** - Proper version tracking and release preparation

## [v0.0.2] - 2025-01-21

## [v0.0.2] - 2025-10-23

### Added
- 🎨 **Beautiful CLI colors and emojis** - Enhanced user experience with colored output
- ✅ **Success messages** - Clear visual feedback for successful operations
- ❌ **Error messages** - Better error reporting with red styling
- 🚀 **Startup messages** - Informative startup information with rocket emoji
- 🌐 **Tunnel status** - Clear indication when tunnels are established
- ℹ️ **Info messages** - Helpful information with blue styling
- ⚠️ **Warning messages** - Important notices with yellow styling

### Improved
- **Better error handling** - More descriptive error messages
- **Enhanced UX** - Visual hierarchy with colors and emojis
- **Cleaner output** - Organized and easy-to-read console output
- **Professional appearance** - Modern CLI that looks polished

### Technical Changes
- Added `internal/colors/colors.go` - Comprehensive color system
- Updated `cmd/client/main.go` - Added colored startup and error messages
- Updated `cmd/server/main.go` - Added colored server startup messages
- Updated `internal/tunnel/client/client.go` - Added colored tunnel status messages
- Removed unused `log` import from server
- Fixed function name collisions in color package

### Bug fixes
- Fix SSL issue with multiple client connected

## [v0.0.1] - 2025-10-21

### Added
- 🚀 **Initial release** - Basic reverse tunnel functionality
- 🌐 **Custom subdomains** - Random subdomain generation
- 🔒 **HTTPS support** - Automatic SSL certificates via Let's Encrypt
- ⚡ **High performance** - Built with Go for speed
- 🔄 **Concurrent tunnels** - Support for multiple simultaneous tunnels
- 🛠️ **Easy setup** - Simple installation and configuration
- 📦 **Cross-platform** - Works on Windows, macOS, and Linux
- 🎯 **Simple CLI** - `got 3000` to tunnel local port 3000
- 🖥️ **Server management** - `server` command for tunnel server
- 🔧 **Health checks** - Automatic tunnel health monitoring
- 📊 **Debug logging** - Comprehensive logging for troubleshooting
- 🎛️ **Configuration** - Environment variables and command-line options
