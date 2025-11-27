package main

import (
	"fmt"
	"github.com/HeyRistaa/got/internal/colors"
)

func main() {
	fmt.Println("=== Mein Tunnel Color Test ===")
	fmt.Println()

	colors.PrintRocket("Starting tunnel for localhost:3000 via 168.119.161.113:4440")
	fmt.Println()

	colors.PrintfSuccess("Tunnel established: %s -> %s\n", colors.Cyan("localhost:3000"), colors.BrightCyan("168.119.161.113:37427"))
	colors.PrintfGlobe("Your service is now available at: %s\n", colors.Bold(colors.BrightGreen("https://abc123.showapps.online")))
	colors.PrintInfo("Press Ctrl+C to stop the tunnel")
	fmt.Println()

	colors.PrintRocket("Starting tunnel server...")
	colors.PrintfInfo("Control port: %s\n", colors.Bold(colors.Blue(":4440")))
	colors.PrintfInfo("Data port: %s\n", colors.Bold(colors.Blue(":4441")))
	colors.PrintSuccess("Server is ready to accept connections!")
	fmt.Println()

	colors.PrintWarning("Health checks disabled")
	colors.PrintError("Connection failed: timeout")
	colors.PrintCheck("Tunnel created successfully")
	colors.PrintCross("Failed to create tunnel")
}
