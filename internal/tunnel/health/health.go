package health

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

// Checker handles health checking for tunnel endpoints
type Checker struct {
	client *http.Client
}

// New creates a new health checker
func New() *Checker {
	return &Checker{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// CheckTunnelEndpoint pings a tunnel endpoint to see if it's healthy
func (c *Checker) CheckTunnelEndpoint(host string) error {
	// Ping the tunnel endpoint via HTTPS
	url := "https://" + host

	resp, err := c.client.Get(url)
	if err != nil {
		// Don't fail on timeout - just log and continue
		if strings.Contains(err.Error(), "timeout") {
			fmt.Printf("Health check timeout for %s (this is normal if client is idle)\n", host)
			return nil // Don't fail on timeout
		}
		return fmt.Errorf("failed to reach tunnel endpoint: %w", err)
	}
	defer resp.Body.Close()

	// Consider 2xx, 3xx, and 4xx as "healthy" (tunnel is working)
	// Only 5xx and connection errors are "unhealthy"
	if resp.StatusCode >= 500 {
		return fmt.Errorf("tunnel endpoint returned %d", resp.StatusCode)
	}

	return nil
}
