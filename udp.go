package main

import (
	"fmt"
	"net"
	"time"
)

// startUDPClient - starts the UDP client and sends a message to the server
func startUDPClient(addr string) {
	var conn net.Conn
	var err error

	//Open up the UDP port and clean up if something goes wrong
	for {
		conn, err = net.Dial("udp", addr)
		if err != nil {
			LogWarningf(true, "could not connect to udp server %v", err)
			time.Sleep((10 * time.Second))
			continue
		}
		break
	}
	defer conn.Close()

	//Once connection using UDP is established, send a message to the server
	// wait random interval between sending messages (MINTIME, MAXTIME)
	for {

		message := fmt.Sprintf("Hello UDP Listener %s", addr)
		conn.Write([]byte(message))

		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err != nil {
			LogWarningf(true, "could not read from UDP server: %v", err)
		}
		responseTxt := string(buffer[:n])
		if len(string(buffer[:n])) == 0 {
			responseTxt = "No response"
		}
		LogInfof(true, "udp connector response from %s: %s\n", addr, responseTxt)
		time.Sleep((randomInterval(MINTIME, MAXTIME)))
	}
}

// startUDPListener - starts up the server on the UDP port.  Sends the response to handleUDPPacket{
func startUDPListener(port string) {

	addr := fmt.Sprintf(":%s", port)
	conn, err := net.ListenPacket("udp", addr)
	if err != nil {
		LogWarningf(true, "starting udp listener on port %s: %v", port, err)
		return
	}
	defer conn.Close()

	LogInfof(true, "udp listener started on port %s", port)

	buffer := make([]byte, 1024)
	for {
		n, clientAddr, err := conn.ReadFrom(buffer)
		if err != nil {
			LogWarningf(true, "could not read from udp packet: %v", err)
			continue
		}
		go handleUDPPacket(conn, clientAddr, buffer[:n])
	}

}

// handleUDPPacket(conn net.PacketConn, clientAddr net.Addr, data []byte) - handles the UDP packet and sends the response back to the client.
func handleUDPPacket(conn net.PacketConn, clientAddr net.Addr, data []byte) {
	currentTime := time.Now().Format(time.RFC3339)
	ip, _, _ := net.SplitHostPort(clientAddr.String())
	response := fmt.Sprintf("Received UDP request from %s at %s\n", ip, currentTime)
	conn.WriteTo([]byte(response), clientAddr)
	LogInfof(true, "udp listener responded with: %s to ip: %s", string(data), ip)
}
