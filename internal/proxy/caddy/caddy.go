package caddy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Client handles Caddy Admin API interactions
type Client struct {
	AdminURL string
}

// New creates a new Caddy client
func New(adminURL string) *Client {
	return &Client{
		AdminURL: adminURL,
	}
}

// AddRoute adds a new route to Caddy
func (c *Client) AddRoute(host string, port int) error {
	fmt.Printf("Adding Caddy route: %s -> 127.0.0.1:%d\n", host, port)

	body := map[string]any{
		"match": []map[string]any{{"host": []string{host}}},
		"handle": []map[string]any{{
			"handler":   "reverse_proxy",
			"upstreams": []map[string]any{{"dial": "127.0.0.1:" + fmt.Sprint(port)}},
			"transport": map[string]any{"protocol": "http", "versions": []string{"1.1"}},
		}},
		// IMPORTANT: Do NOT set terminal=true so that ACME HTTP-01 challenge handlers
		// can intercept /.well-known/acme-challenge/* before our reverse_proxy route.
	}
	var buf bytes.Buffer
	_ = json.NewEncoder(&buf).Encode(body)

	req, err := http.NewRequest("POST", c.AdminURL+"/config/apps/http/servers/srv0/routes", &buf)
	if err != nil {
		fmt.Printf("Failed to create Caddy request for %s: %v\n", host, err)
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("Failed to send Caddy request for %s: %v\n", host, err)
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		fmt.Printf("Caddy route add failed for %s: %s - %s\n", host, resp.Status, strings.TrimSpace(string(b)))
		return fmt.Errorf("caddy add route: %s: %s", resp.Status, strings.TrimSpace(string(b)))
	}
	fmt.Printf("Successfully added Caddy route for %s\n", host)
	return nil
}

// DeleteRouteByHost removes a route by hostname
func (c *Client) DeleteRouteByHost(host string) error {
	// 1) GET routes to find the index containing this host
	resp, err := http.Get(c.AdminURL + "/config/apps/http/servers/srv0/routes")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var routes []map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&routes); err != nil {
		return err
	}

	var idx = -1
	for i, r := range routes {
		ms, _ := r["match"].([]any)
		for _, m := range ms {
			if mm, ok := m.(map[string]any); ok {
				if hosts, ok := mm["host"].([]any); ok {
					for _, h := range hosts {
						if hs, ok := h.(string); ok && hs == host {
							idx = i
							break
						}
					}
				}
			}
		}
		if idx >= 0 {
			break
		}
	}
	if idx < 0 {
		return nil // nothing to delete
	}

	// 2) DELETE that route index
	req, err := http.NewRequest("DELETE",
		fmt.Sprintf("%s/config/apps/http/servers/srv0/routes/%d", c.AdminURL, idx), nil)
	if err != nil {
		return err
	}
	resp2, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp2.Body.Close()
	if resp2.StatusCode < 200 || resp2.StatusCode >= 300 {
		b, _ := io.ReadAll(resp2.Body)
		return fmt.Errorf("caddy del route: %s: %s", resp2.Status, strings.TrimSpace(string(b)))
	}
	return nil
}
