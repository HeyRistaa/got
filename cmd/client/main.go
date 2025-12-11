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

	"github.com/HeyRistaa/got/internal/colors"
	"github.com/HeyRistaa/got/internal/tunnel/client"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func main() {
	var server string
	var local string
	var id string
	var domain string
	flag.StringVar(&server, "server", "", "server host (port will be 4440)")
	flag.StringVar(&local, "local", "", "local address to forward")
	flag.StringVar(&id, "id", "", "client identifier")
	flag.StringVar(&domain, "domain", "", "domain to use for the tunnel")
	flag.Parse()

	// Get server host - priority: CLI > env > Hetzner API
	// Security: Do NOT hardcode production server IP in public repo
	serverHost := ""
	if server != "" {
		serverHost = server
	} else if envHost := os.Getenv("GOT_SERVER_HOST"); envHost != "" {
		serverHost = envHost
	} else if ip := resolveHetznerIPFromEnv(); ip != "" {
		serverHost = ip
	} else {
		colors.PrintError("Error: Server not specified. Please provide a server host.\n")
		colors.PrintInfo("Usage: got -server <host> <port>\n")
		colors.PrintInfo("Or set GOT_SERVER_HOST environment variable\n")
		os.Exit(1)
	}

	// Fixed ports - no environment variables needed
	controlAddr := fmt.Sprintf("%s:4440", serverHost)
	dataAddr := fmt.Sprintf("%s:4441", serverHost)

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
		colors.PrintError("Usage: got -server <host> <localPort|host:port> [flags]\n")
		colors.PrintInfo("Example: got -server your-server.com 3000\n")
		colors.PrintInfo("Or set GOT_SERVER_HOST environment variable and use: got 3000\n")
		os.Exit(1)
	}

	colors.PrintRocket("Starting tunnel for " + colors.Cyan(local) + " via " + colors.Blue(controlAddr) + "\n")

	c := client.New(controlAddr, dataAddr, local, id, domain)
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	if err := c.Run(ctx); err != nil {
		colors.PrintfError("Client error: %v\n", err)
		os.Exit(1)
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
