package control

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
)

// Message types exchanged over the persistent control connection and
// short-lived data init connections. Newline-delimited JSON frames.

// Server to open a tunnel
// Client -> Server
type OpenTunnel struct {
	Type      string `json:"type"`       // "open_tunnel"
	ClientID  string `json:"client_id"`  // optional identifier
	LocalHint string `json:"local_hint"` // optional label for debugging
	Domain    string `json:"domain"`     // optional domain to use for the tunnel
	LocalURL  string `json:"local_url"`  // local URL for health checking (e.g., "http://localhost:3000")
}

// Server to client with a TunnelID and the public host:port on the server accessible publicly
type TunnelOpened struct {
	Type       string `json:"type"` // "tunnel_opened"
	TunnelID   string `json:"tunnel_id"`
	PublicAddr string `json:"public_addr"` // host:port on server accessible publicly
	PublicHost string `json:"public_host"` // optional host when using host-based routing
}

// Server to client with an error message
type TunnelError struct {
	Type  string `json:"type"` // "tunnel_error"
	Error string `json:"error"`
}

// Client asking it to open a data connection to the server for incoming connections
type ConnRequest struct {
	Type     string `json:"type"` // "conn_request"
	TunnelID string `json:"tunnel_id"`
	ConnID   string `json:"conn_id"`
}

// Sent by client on a short-lived TCP connection to initialize a data pipe.
type DataInit struct {
	Type     string `json:"type"` // "data_init"
	TunnelID string `json:"tunnel_id"`
	ConnID   string `json:"conn_id"`
}

// Heartbeat message from client to server
type Heartbeat struct {
	Type     string `json:"type"` // "heartbeat"
	TunnelID string `json:"tunnel_id"`
}

// JSON line helpers

func WriteJSONLine(w io.Writer, v any) error {
	enc := json.NewEncoder(w)
	if err := enc.Encode(v); err != nil {
		return err
	}
	return nil
}

func ReadJSONLine(r *bufio.Reader, v any) error {
	line, err := r.ReadBytes('\n')
	if err != nil {
		return err
	}
	if err := json.Unmarshal(line, v); err != nil {
		return fmt.Errorf("invalid json: %w", err)
	}
	return nil
}
