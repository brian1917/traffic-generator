package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

// startTCPListener - starts up the server on the TCP port.  Sends the response to handlerTCPConnection
func startTCPListener(port string) {
	addr := fmt.Sprintf(":%s", port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Error starting TCP listener on port %s: %v", port, err)
	}
	defer listener.Close()

	log.Printf("TCP listener started on port %s", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting TCP connection: %v", err)
			continue
		}
		go handleTCPConnection(conn)
	}
}

// handleTCPConnection(conn net.Conn) - handles the TCP connection and sends the response back to the client.
func handleTCPConnection(conn net.Conn) {
	defer conn.Close()
	ip, _, _ := net.SplitHostPort(conn.RemoteAddr().String())
	currentTime := time.Now().Format(time.RFC3339)
	response := fmt.Sprintf("Received TCP request from %s at %s\n", ip, currentTime)
	conn.Write([]byte(response))
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		log.Printf("Error reading TCP packet: %v", err)
	}
	log.Printf("TCP Listener responded to with: %s to ip: %s", string(buffer[:n]), ip)
}

// startTCPClient - starts the TCP client and sends a message to the server
func startTCPClient(addr string) {
	var conn net.Conn
	var err error

	for {
		// Try to connect to the server
		conn, err = net.Dial("tcp", addr)
		if err != nil {
			log.Printf("Error connecting to TCP server %s: %v", addr, err)
			time.Sleep(10 * time.Second) // Retry after 10 seconds if connection fails
			continue
		}

		// Ensure the connection is closed when done
		defer conn.Close()

		// Set a read deadline to avoid blocking forever
		conn.SetReadDeadline(time.Now().Add(10 * time.Second))

		// Send a message to the server
		message := fmt.Sprintf("Hello TCP Listener %s", addr)
		_, err = conn.Write([]byte(message))
		if err != nil {
			log.Printf("Error writing to TCP server %s: %v", addr, err)
			conn.Close()
			time.Sleep(randomInterval(3, 5)) // Sleep before retrying
			continue
		}

		// Read the server's response
		buffer := make([]byte, 1024)

		// Keep reading until we get a response
		for {
			n, err := conn.Read(buffer)
			if err != nil {
				// Check for specific errors
				if err.Error() == "i/o timeout" {
					log.Println("Read timed out, retrying...")
				} else {
					log.Printf("Error reading from TCP server %s: %v", addr, err)
				}
				// Sleep and try again
				time.Sleep(randomInterval(3, 5)) // Random sleep before retry
				// Reset the read deadline after sleeping
				conn.SetReadDeadline(time.Now().Add(10 * time.Second))
				continue
			}

			// If we get data, log it and break out of the read loop
			log.Printf("TCP Connector response from %s: %s\n", addr, string(buffer[:n]))
			break
		}

		// wait random interval between sending messages (MINTIME, MAXTIME)
		time.Sleep(randomInterval(MINTIME, MAXTIME))
	}
}
