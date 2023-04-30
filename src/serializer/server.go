package main

import (
	"fmt"
	"net"
	"strings"
)

func main() {
	// Calculate the result in a separate goroutine
	result := make(chan string)
	go func() {
		result <- calculateResult()
	}()

	// Start the server
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(fmt.Sprintf("Failed to start server: %v", err))
	}
	defer listener.Close()
	fmt.Println("Server started")

	var r string

	// Listen for incoming requests
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Failed to accept connection: %v\n", err)
			continue
		}

		// Read incoming message
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Printf("Failed to read message: %v\n", err)
			conn.Close()
			continue
		}
		msg := string(buf[:n])

		// Handle message
		if strings.TrimSpace(msg) == "get_result" {
			if r == "" {
				select {
				case r = <-result:
					// Do nothing, will write r later
				default:
					conn.Write([]byte("not ready\n"))
				}
			}

			if r != "" {
				conn.Write([]byte(r))
			}

		} else {
			conn.Write([]byte(fmt.Sprintf("unknown command: %s\n", msg)))
		}

		// Close connection
		conn.Close()
	}
}
