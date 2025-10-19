package tunnel

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/HeyRistaa/got/internal/proxy/caddy"
	"github.com/HeyRistaa/got/internal/tunnel/health"
)

// Manager handles tunnel lifecycle
type Manager struct {
	caddyClient *caddy.Client
	healthChecker *health.Checker
}

// Tunnel represents a single tunnel
type Tunnel struct {
	ID       string
	ClientID string
	Port     int
	Host     string
	Domain   string
	Listener net.Listener
}

// NewManager creates a new tunnel manager
func NewManager() *Manager {
	return &Manager{
		caddyClient:   caddy.New("http://127.0.0.1:2019"),
		healthChecker: health.New(),
	}
}

// CreateTunnel creates a new tunnel
func (m *Manager) CreateTunnel(clientID, domain string) (*Tunnel, error) {
	// Generate random subdomain
	label := randomID()[:6]
	if strings.HasPrefix(domain, "*.") {
		domain = strings.TrimPrefix(domain, "*.")
	}
	host := fmt.Sprintf("%s.%s", label, domain)

	// Allocate port
	port, listener, err := m.allocatePort()
	if err != nil {
		return nil, fmt.Errorf("failed to allocate port: %w", err)
	}

	// Create Caddy route
	if err := m.caddyClient.AddRoute(host, port); err != nil {
		listener.Close()
		return nil, fmt.Errorf("failed to add caddy route: %w", err)
	}

	tunnel := &Tunnel{
		ID:       randomID(),
		ClientID: clientID,
		Port:     port,
		Host:     host,
		Domain:   domain,
		Listener: listener,
	}

	return tunnel, nil
}

// StartHealthCheck starts health checking for a tunnel
func (m *Manager) StartHealthCheck(tunnel *Tunnel, cleanupFunc func()) {
	go func() {
		ticker := time.NewTicker(30 * time.Second) // Check every 30 seconds
		defer ticker.Stop()
		
		for {
			select {
			case <-ticker.C:
				if err := m.healthChecker.CheckTunnelEndpoint(tunnel.Host); err != nil {
					fmt.Printf("Tunnel endpoint health check failed for %s: %v\n", tunnel.Host, err)
					cleanupFunc()
					return
				}
				fmt.Printf("Tunnel endpoint health check passed for %s\n", tunnel.Host)
			}
		}
	}()
}

// CloseTunnel closes a tunnel and cleans up resources
func (m *Manager) CloseTunnel(tunnel *Tunnel) error {
	if tunnel.Listener != nil {
		tunnel.Listener.Close()
	}
	
	// Remove Caddy route
	if err := m.caddyClient.DeleteRouteByHost(tunnel.Host); err != nil {
		return fmt.Errorf("failed to delete caddy route: %w", err)
	}
	
	return nil
}

// allocatePort allocates a port for the tunnel
func (m *Manager) allocatePort() (int, net.Listener, error) {
	// If PUBLIC_PORT env is set, try to bind that exact port; otherwise use :0
	if forced := os.Getenv("PUBLIC_PORT"); forced != "" {
		if _, err := strconv.Atoi(forced); err == nil {
			if ln, err := net.Listen("tcp", net.JoinHostPort("", forced)); err == nil {
				port, _ := strconv.Atoi(forced)
				return port, ln, nil
			} else {
				return 0, nil, fmt.Errorf("failed to bind PUBLIC_PORT %s: %w", forced, err)
			}
		}
	}
	// Bind to :0 to get a free port
	ln, err := net.Listen("tcp", net.JoinHostPort("", "0"))
	if err != nil {
		return 0, nil, err
	}
	addr := ln.Addr().(*net.TCPAddr)
	port := addr.Port
	// Keep this listener to accept users directly
	return port, ln, nil
}

// randomID generates a random ID
func randomID() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
