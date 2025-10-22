package server

import (
	"bufio"
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net"
	"net/netip"
	"sync"
	"time"

	"github.com/HeyRistaa/got/internal/protocol/control"
	"github.com/HeyRistaa/got/internal/tunnel"
)

type Server struct {
	ControlListen string // host:port where clients establish control connections
	DataListen    string // host:port where SERVER accepts client data connections
	PublicIP      string // public IP or hostname to advertise (e.g., Hetzner IP)

	mu      sync.RWMutex
	tunnels map[string]*tunnelInfo // by tunnelID
	ports   map[int]string         // public port -> tunnelID

	pendingMu sync.Mutex
	pending   map[string]chan net.Conn // connID -> ready data conn from client

	tunnelManager *tunnel.Manager
}

type tunnelInfo struct {
	tunnel  *tunnel.Tunnel
	ctlConn net.Conn
}

func New(controlAddr, dataAddr, publicIP string) *Server {
	return &Server{
		ControlListen: controlAddr,
		DataListen:    dataAddr,
		PublicIP:      publicIP,
		tunnels:       make(map[string]*tunnelInfo),
		ports:         make(map[int]string),
		pending:       make(map[string]chan net.Conn),
		tunnelManager: tunnel.NewManager(),
	}
}

func (s *Server) Run(ctx context.Context) error {
	// listener for control connections (from clients)
	ctlLn, err := net.Listen("tcp", s.ControlListen)
	if err != nil {
		return fmt.Errorf("listen control: %w", err)
	}
	defer ctlLn.Close()

	// listener for client data connections
	dataLn, err := net.Listen("tcp", s.DataListen)
	if err != nil {
		return fmt.Errorf("listen data: %w", err)
	}
	defer dataLn.Close()

	log.Printf("server: control %s, data %s, public IP %s", s.ControlListen, s.DataListen, s.PublicIP)

	// Accept control connections and handle in goroutines
	go func() {
		for {
			conn, err := ctlLn.Accept()
			if err != nil {
				if errors.Is(err, net.ErrClosed) {
					return
				}
				log.Printf("accept control: %v", err)
				continue
			}
			go s.handleControl(conn)
		}
	}()

	// Accept client data connections and match by DataInit
	go func() {
		for {
			conn, err := dataLn.Accept()
			if err != nil {
				if errors.Is(err, net.ErrClosed) {
					return
				}
				log.Printf("accept data: %v", err)
				continue
			}
			go s.handleDataConn(conn)
		}
	}()

	// Block until context done
	<-ctx.Done()
	return nil
}

func (s *Server) handleControl(conn net.Conn) {
	defer conn.Close()
	r := bufio.NewReader(conn)

	var req control.OpenTunnel
	if err := control.ReadJSONLine(r, &req); err != nil {
		log.Printf("control: read open_tunnel: %v", err)
		return
	}
	if req.Type != "open_tunnel" {
		_ = control.WriteJSONLine(conn, control.TunnelError{Type: "tunnel_error", Error: "expected open_tunnel"})
		return
	}

	log.Printf("Received open_tunnel request from client %s for local %s", req.ClientID, req.LocalHint)
	log.Printf("Current active tunnels: %d", len(s.tunnels))

	// Create tunnel using tunnel manager
	domain := req.Domain
	if domain == "" {
		domain = "*.showapps.online" // default
	}

	log.Printf("Creating tunnel for client %s with domain %s", req.ClientID, domain)
	tunnel, err := s.tunnelManager.CreateTunnel(req.ClientID, domain)
	if err != nil {
		log.Printf("Failed to create tunnel for client %s: %v", req.ClientID, err)
		_ = control.WriteJSONLine(conn, control.TunnelError{Type: "tunnel_error", Error: err.Error()})
		return
	}
	log.Printf("Successfully created tunnel %s for client %s", tunnel.ID, req.ClientID)

	// Store tunnel info
	tid := tunnel.ID
	s.mu.Lock()
	s.tunnels[tid] = &tunnelInfo{
		tunnel:  tunnel,
		ctlConn: conn,
	}
	s.ports[tunnel.Port] = tid
	s.mu.Unlock()

	log.Printf("Stored tunnel %s (port %d) for client %s", tid, tunnel.Port, req.ClientID)

	// Send tunnel opened response
	pub := fmt.Sprintf("%s:%d", s.PublicIP, tunnel.Port)
	opened := control.TunnelOpened{
		Type:       "tunnel_opened",
		TunnelID:   tid,
		PublicAddr: pub,
		PublicHost: tunnel.Host,
	}
	if err := control.WriteJSONLine(conn, opened); err != nil {
		log.Printf("control: write opened: %v", err)
		return
	}

	// Start serving public connections
	log.Printf("Starting public listener for tunnel %s on port %d", tid, tunnel.Port)
	go s.servePublic(tunnel.Listener, tid, tunnel.Port)

	// Start health checking
	log.Printf("Starting health check for tunnel %s", tid)
	s.tunnelManager.StartHealthCheck(tunnel, func() {
		log.Printf("Health check cleanup triggered for tunnel %s", tid)
		s.cleanupTunnel(tid, tunnel.Port, tunnel.Listener)
	})

	// Keep control connection open until client disconnects
	log.Printf("Control connection established for tunnel %s, waiting for client disconnect", tid)
	_ = conn.SetReadDeadline(time.Time{})
	for {
		if _, err := r.Peek(1); err != nil {
			// client closed
			log.Printf("Client disconnected for tunnel %s, cleaning up", tid)
			s.cleanupTunnel(tid, tunnel.Port, tunnel.Listener)
			return
		}
		time.Sleep(5 * time.Second)
	}
}

func (s *Server) servePublic(ln net.Listener, tunnelID string, port int) {
	defer ln.Close()
	for {
		userConn, err := ln.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return
			}
			log.Printf("public accept [%d]: %v", port, err)
			continue
		}
		go s.bridgeUserConnection(tunnelID, userConn)
	}
}

func (s *Server) bridgeUserConnection(tunnelID string, userConn net.Conn) {
	// Ask client to open a data connection back to server
	s.mu.RLock()
	t := s.tunnels[tunnelID]
	s.mu.RUnlock()
	if t == nil || t.ctlConn == nil {
		log.Printf("no control conn for tunnel %s", tunnelID)
		userConn.Close()
		return
	}

	connID := randomID()
	ch := make(chan net.Conn, 1)
	s.pendingMu.Lock()
	s.pending[connID] = ch
	s.pendingMu.Unlock()

	if err := control.WriteJSONLine(t.ctlConn, control.ConnRequest{Type: "conn_request", TunnelID: tunnelID, ConnID: connID}); err != nil {
		log.Printf("write conn_request: %v", err)
		userConn.Close()
		s.clearPending(connID)
		return
	}

	select {
	case dataConn := <-ch:
		pipeConns(userConn, dataConn)
	case <-time.After(10 * time.Second):
		log.Printf("timeout waiting for client data conn for %s", connID)
		userConn.Close()
	}
	s.clearPending(connID)
}

func (s *Server) cleanupTunnel(tunnelID string, port int, ln net.Listener) {
	log.Printf("Cleaning up tunnel %s (port %d)", tunnelID, port)
	ln.Close()
	s.mu.Lock()
	t := s.tunnels[tunnelID]
	delete(s.tunnels, tunnelID)
	delete(s.ports, port)
	s.mu.Unlock()

	// Close tunnel using tunnel manager
	if t != nil && t.tunnel != nil {
		log.Printf("Closing tunnel %s via tunnel manager", tunnelID)
		if err := s.tunnelManager.CloseTunnel(t.tunnel); err != nil {
			log.Printf("failed to close tunnel: %v", err)
		}
	}
	log.Printf("Tunnel %s cleanup completed", tunnelID)
}

func (s *Server) handleDataConn(conn net.Conn) {
	r := bufio.NewReader(conn)
	var init control.DataInit
	if err := control.ReadJSONLine(r, &init); err != nil {
		log.Printf("data init read: %v", err)
		conn.Close()
		return
	}
	if init.Type != "data_init" {
		conn.Close()
		return
	}
	s.pendingMu.Lock()
	ch := s.pending[init.ConnID]
	s.pendingMu.Unlock()
	if ch == nil {
		conn.Close()
		return
	}
	ch <- conn
}

func (s *Server) clearPending(connID string) {
	s.pendingMu.Lock()
	delete(s.pending, connID)
	s.pendingMu.Unlock()
}

// Helper functions (keeping the existing ones)
func randomID() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func pipeConns(a, b net.Conn) {
	var wg sync.WaitGroup
	copy := func(dst, src net.Conn) {
		defer wg.Done()
		_, _ = ioCopy(dst, src)
		dst.Close()
	}
	wg.Add(2)
	go copy(a, b)
	go copy(b, a)
	wg.Wait()
}

// ioCopy is a thin wrapper to allow deadline tweaks later.
func ioCopy(dst net.Conn, src net.Conn) (int64, error) {
	return netCopy(dst, src)
}

// netCopy mirrors io.Copy but keeps types explicit for future tuning.
func netCopy(dst net.Conn, src net.Conn) (written int64, err error) {
	buf := make([]byte, 32*1024)
	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			nw, ew := dst.Write(buf[0:nr])
			if nw > 0 {
				written += int64(nw)
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = fmt.Errorf("short write")
				break
			}
		}
		if er != nil {
			if errors.Is(er, net.ErrClosed) {
				return written, nil
			}
			return written, er
		}
	}
	return written, err
}

// Helper to validate IP for logs
func isValidHostPort(addr string) bool {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return false
	}
	_, err = netip.ParseAddr(host)
	return err == nil
}
