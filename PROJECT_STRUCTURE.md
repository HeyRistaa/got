# Go Project Structure - Standard Layout

## 📁 New Project Structure

```
got/
├── cmd/                          # Main applications
│   ├── client/                   # Client CLI application
│   │   └── main.go              # Entry point for 'got' command
│   └── server/                   # Server application  
│       └── main.go              # Entry point for 'server' command
│
├── internal/                     # Private application code
│   ├── protocol/                 # Communication protocols
│   │   └── control/             # Control protocol definitions
│   │       └── protocol.go      # JSON message types
│   │
│   ├── proxy/                    # Proxy-related functionality
│   │   └── caddy/               # Caddy integration
│   │       └── caddy.go         # Caddy Admin API client
│   │
│   └── tunnel/                   # Core tunnel functionality
│       ├── client/              # Client implementation
│       │   └── client.go       # Client logic
│       ├── server/              # Server implementation
│       │   └── server.go       # Server logic
│       ├── health/              # Health checking
│       │   └── health.go       # Health check utilities
│       └── manager.go           # Tunnel lifecycle management
│
├── releases/                     # Pre-built binaries
├── install.sh                   # Installation script
├── install.ps1                  # Windows installation script
├── go.mod                       # Go module definition
├── LICENSE                      # MIT License
└── README.md                    # Project documentation
```

## 🎯 Benefits of This Structure

### **1. Standard Go Layout Compliance**
- Follows [golang-standards/project-layout](https://github.com/golang-standards/project-layout)
- Industry-recognized best practice
- Easy for other Go developers to understand

### **2. Clear Separation of Concerns**
- **`cmd/`**: Entry points only, minimal logic
- **`internal/protocol/`**: Communication protocols
- **`internal/proxy/`**: External service integrations (Caddy)
- **`internal/tunnel/`**: Core business logic

### **3. Logical Grouping**
- **Protocol**: All communication-related code
- **Proxy**: External service integrations
- **Tunnel**: Core tunnel functionality (client, server, health, manager)

### **4. Scalability**
- Easy to add new protocols (WebSocket, gRPC, etc.)
- Easy to add new proxy integrations (Nginx, Traefik, etc.)
- Easy to extend tunnel functionality

### **5. Maintainability**
- Clear package boundaries
- Easy to locate specific functionality
- Reduced coupling between components

## 🔄 Migration Summary

### **Before (Messy Structure):**
```
internal/
├── caddy/          # Mixed with other concerns
├── health/         # Scattered
├── tunnel/         # Confusing naming
├── control/        # Generic name
├── client/         # Mixed with server
└── server/         # Mixed with client
```

### **After (Clean Structure):**
```
internal/
├── protocol/control/    # Clear protocol focus
├── proxy/caddy/         # Clear proxy focus  
└── tunnel/              # All tunnel-related code
    ├── client/          # Client implementation
    ├── server/          # Server implementation
    ├── health/          # Health checking
    └── manager.go       # Lifecycle management
```

## 🚀 Usage

### **Build Commands:**
```bash
# Build client
go build -o got ./cmd/client

# Build server  
go build -o server ./cmd/server

# Build both
go build -o got ./cmd/client && go build -o server ./cmd/server
```

### **Import Examples:**
```go
// Client imports tunnel client
import "github.com/HeyRistaa/got/internal/tunnel/client"

// Server imports tunnel server
import "github.com/HeyRistaa/got/internal/tunnel/server"

// Tunnel manager imports proxy and health
import "github.com/HeyRistaa/got/internal/proxy/caddy"
import "github.com/HeyRistaa/got/internal/tunnel/health"
```

This structure is now **professional**, **scalable**, and follows **Go best practices**! 🎉
