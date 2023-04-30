package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

func main() {
	proxyAddress := ":2000"
	proxyListener, err := net.Listen("tcp", proxyAddress)
	if err != nil {
		log.Fatalf("Failed to start proxy: %s", err)
		return
	}
	defer proxyListener.Close()

	log.Printf("Proxy is running on %s", proxyAddress)

	for {
		conn, err := proxyListener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %s", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				log.Printf("Error reading line: %s", err)
			}
			return
		}

		line = strings.TrimSpace(line)
		parts := strings.Split(line, " ")
		if len(parts) != 2 || parts[0] != "get_result" {
			conn.Write([]byte("Invalid message format\n"))
			continue
		}

		format := parts[1]
		backendAddress := resolveHostForFormat(format)
		if backendAddress == "" {
			conn.Write([]byte("Invalid format\n"))
			continue
		}

		result, err := forwardRequest(backendAddress)
		if err != nil {
			conn.Write([]byte(fmt.Sprintf("Error forwarding request: %s\n", err)))
			continue
		}

		conn.Write([]byte(result))
	}
}

func resolveHostForFormat(format string) string {
	switch format {
	case "gob":
		return "serializer-gob:8080"
	case "xml":
		return "serializer-xml:8080"
	case "json":
		return "serializer-json:8080"
	case "protobuf":
		return "serializer-protobuf:8080"
	case "avro":
		return "serializer-avro:8080"
	case "yaml":
		return "serializer-yaml:8080"
	case "message-pack":
		return "serializer-message-pack:8080"
	default:
		return ""
	}
}

func forwardRequest(backendAddress string) (string, error) {
	conn, err := net.Dial("tcp", backendAddress)
	if err != nil {
		return "", fmt.Errorf("Failed to connect to backend: %s", err)
	}
	defer conn.Close()

	_, err = conn.Write([]byte("get_result\n"))
	if err != nil {
		return "", fmt.Errorf("Failed to send get_result to backend: %s", err)
	}

	reader := bufio.NewReader(conn)
	result, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("Failed to read response from backend: %s", err)
	}

	return result, nil
}
