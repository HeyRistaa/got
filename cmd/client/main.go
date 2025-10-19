package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/HeyRistaa/got/internal/tunnel/client"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func main() {
	var server string
	var data string
	var local string
	var id string
	var domain string
	flag.StringVar(&server, "server", "", "server control address (host or host:port)")
	flag.StringVar(&data, "data", "", "server data address (host:port)")
	flag.StringVar(&local, "local", "", "local address to forward")
	flag.StringVar(&id, "id", "", "client identifier")
	flag.StringVar(&domain, "domain", "", "domain to use for the tunnel")
	flag.Parse()

	// Defaults from environment or sensible fallbacks
	serverHost := os.Getenv("GOT_SERVER_HOST")
	if serverHost == "" {
		if ip := resolveHetznerIPFromEnv(); ip != "" {
			serverHost = ip
		} else {
			serverHost = "127.0.0.1"
		}
	}
	controlPort := os.Getenv("GOT_CONTROL_PORT")
	if controlPort == "" {
		controlPort = "4440"
	}
	dataPort := os.Getenv("GOT_DATA_PORT")
	if dataPort == "" {
		dataPort = "4441"
	}

	// If --server provided, it may be host or host:port
	if server != "" {
		if strings.Contains(server, ":") {
			// Keep as-is for control, but extract host for data default
			serverHost = strings.Split(server, ":")[0]
			// respect provided full control address
		} else {
			serverHost = server
			server = fmt.Sprintf("%s:%s", serverHost, controlPort)
		}
	} else {
		server = fmt.Sprintf("%s:%s", serverHost, controlPort)
	}

	// If --data not provided, derive from server host and dataPort
	if data == "" {
		data = fmt.Sprintf("%s:%s", serverHost, dataPort)
	}

	// Positional arg convenience: `got 3002` or `got localhost:3002`
	if local == "" && flag.NArg() >= 1 {
		arg := flag.Arg(0)
		if strings.Contains(arg, ":") {
			local = arg
		} else if _, err := strconv.Atoi(arg); err == nil {
			// numeric port
			local = fmt.Sprintf("localhost:%s", arg)
		} else {
			// treat as host (missing port)
			log.Fatalf("invalid local argument: %s (expected port or host:port)", arg)
		}
	}

	if local == "" {
		log.Fatalf("usage: got <localPort|host:port> [flags]. Example: got 3002 or got -domain '*.apps.mydomain.com' 3002")
	}

	c := client.New(server, data, local, id, domain)
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	if err := c.Run(ctx); err != nil {
		log.Fatalf("client error: %v", err)
	}
}

// resolveHetznerIPFromEnv tries to read GOT_HC_TOKEN and either GOT_HC_SERVER_ID
// or GOT_HC_SERVER_NAME to obtain the VM public IPv4 via Hetzner Cloud API.
// Returns empty string on any failure or if envs are not set.
func resolveHetznerIPFromEnv() string {
	token := os.Getenv("GOT_HC_TOKEN")
	if token == "" {
		return ""
	}
	client := hcloud.NewClient(hcloud.WithToken(token))
	if idStr := os.Getenv("GOT_HC_SERVER_ID"); idStr != "" {
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return ""
		}
		if srv, _, err := client.Server.GetByID(context.Background(), id); err == nil && srv != nil {
			if ip := srv.PublicNet.IPv4.IP; ip != nil {
				return ip.String()
			}
		}
		return ""
	}
	if name := os.Getenv("GOT_HC_SERVER_NAME"); name != "" {
		if srv, _, err := client.Server.GetByName(context.Background(), name); err == nil && srv != nil {
			if ip := srv.PublicNet.IPv4.IP; ip != nil {
				return ip.String()
			}
		}
	}
	return ""
}
