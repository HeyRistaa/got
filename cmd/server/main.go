package main

import (
	"context"
	"flag"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/HeyRistaa/got/internal/colors"
	"github.com/HeyRistaa/got/internal/tunnel/server"
)

func main() {
	var publicIP string
	var disableHealthCheck bool
	flag.StringVar(&publicIP, "public", "", "public IP/host advertised for tunnels")
	flag.BoolVar(&disableHealthCheck, "disable-health-check", false, "disable health checks for tunnels")
	flag.Parse()

	if publicIP == "" {
		colors.PrintInfo("Detecting public IP...\n")
		publicIP = detectPublicIP()
	}
	if publicIP == "" {
		colors.PrintError("Could not detect public IP, please provide it with the -public flag\n")
		os.Exit(1)
	}
	colors.PrintfInfo("Public IP: %s\n", colors.Bold(colors.BrightCyan(publicIP)))

	// Set environment variable for health check disable
	if disableHealthCheck {
		os.Setenv("GOT_DISABLE_HEALTH_CHECK", "1")
		colors.PrintWarning("Health checks disabled\n")
	}

	// Fixed ports for simplicity
	controlAddr := ":4440"
	dataAddr := ":4441"

	colors.PrintRocket("Starting tunnel server...\n")
	colors.PrintfInfo("Control port: %s\n", colors.Bold(colors.Blue(controlAddr)))
	colors.PrintfInfo("Data port: %s\n", colors.Bold(colors.Blue(dataAddr)))
	colors.PrintSuccess("Server is ready to accept connections!\n")

	srv := server.New(controlAddr, dataAddr, publicIP)
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	if err := srv.Run(ctx); err != nil {
		colors.PrintfError("Server error: %v\n", err)
		os.Exit(1)
	}
}

func detectPublicIP() string {
	resp, err := http.Get("https://api.ipify.org")
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(body))
}
