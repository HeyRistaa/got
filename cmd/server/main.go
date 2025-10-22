package main

import (
	"context"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/HeyRistaa/got/internal/tunnel/server"
)

func main() {
	var publicIP string
	var disableHealthCheck bool
	flag.StringVar(&publicIP, "public", "", "public IP/host advertised for tunnels")
	flag.BoolVar(&disableHealthCheck, "disable-health-check", false, "disable health checks for tunnels")
	flag.Parse()

	if publicIP == "" {
		publicIP = detectPublicIP()
	}
	if publicIP == "" {
		log.Fatalf("Could not detect public IP, please provide it with the -public flag")
	}
	log.Printf("Public IP: %s", publicIP)

	// Set environment variable for health check disable
	if disableHealthCheck {
		os.Setenv("GOT_DISABLE_HEALTH_CHECK", "1")
		log.Printf("Health checks disabled")
	}

	// Fixed ports for simplicity
	controlAddr := ":4440"
	dataAddr := ":4441"
	srv := server.New(controlAddr, dataAddr, publicIP)
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	if err := srv.Run(ctx); err != nil {
		log.Fatalf("server error: %v", err)
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
