package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"
)

// clientTraffic - struct to hold the ports that need to be open for inbound and the outbound client traffic data
type clientTraffic struct {
	listenerTCPPorts []string
	listenerUDPPorts []string
	dstTCPconnection []string
	dstUDPconnection []string
}

// MINTIME and MAXTIME are the minimum and maximum time in seconds to wait between sending traffic
const MINTIME = 60
const MAXTIME = 600

func main() {

	// Check error
	if len(os.Args) < 2 {
		log.Fatal("at least one argument needed: open or traffic")
	}

	if os.Args[1] == "open" && len(os.Args) != 3 {
		log.Fatal("open requires another argument for port (e.g., traffic-generator open 3306)")
	}

	if os.Args[1] == "traffic" && len(os.Args) != 3 {
		log.Fatal("traffic requires another argument for the traffic import file. the import file should have 3 headers: src, dst, and port.")
	}

	if os.Args[1] == "vensim" && len(os.Args) != 3 {
		log.Fatal("vensim requires another argument for the vensim import file. the import file should have 14 headers: src_ip,src_hostname,src_process,src_service,src_username,dst_ip,dst_hostname,dst_fqdn,port,proto,dst_process,dst_service,dst_username,policy_decisiont.")
	}

	// Call open
	if os.Args[1] == "open" {
		port, err := strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatalf("cannot process port - %s", err)

		}
		openListeningPort(port)
	}

	// Call traffic reads a CSV file and sends traffic to the appropriate ports
	if os.Args[1] == "traffic" {
		// Validate hostname
		log.Printf("hostname discovered as %s\r\n", hostname())
		sendTraffic(os.Args[2])
	}

	// Call vensim reads vensim traffic CSV file.  Runs a listener for ports the ip matches on and sends traffic where it is the sourc
	if os.Args[1] == "vensim" {
		// Validate hostname
		log.Printf("hostname discovered as %s\r\n", hostname())
		vensimTraffic(os.Args[2])
	}

}

/*
* Headers
src - hostname srouce
dst - hostname destination
port - destination port
*
*/
func sendTraffic(csvFile string) {
	// Open CSV File and create the reader
	file, err := os.Open(csvFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	reader := csv.NewReader(bufio.NewReader(file))
	lineNumber := 0

	// Iterate through the CSV
	for {
		// Read the line
		line, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		// Increment the line number and skip header row
		lineNumber++
		if lineNumber == 1 {
			continue
		}

		// Skip if it's not for this current host
		if line[0] != hostname() {
			continue
		}

		// Make HTTP request for
		resp, err := http.Get(fmt.Sprintf("http://%s:%s", line[1], line[2]))
		if err != nil {
			log.Fatal(err)
		}
		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("error reading response body - %s", err)
			return
		}
		fmt.Println(string(responseBody))
	}
}

func LoadCSV(csvFile string) ([][]string, map[string]int, error) {
	// Open CSV File and create the reader
	file, err := os.Open(csvFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	reader := csv.NewReader(ClearBOM(bufio.NewReader(file)))
	lineNumber := 0

	var data [][]string
	// Iterate through the CSV

	headerMap := make(map[string]int)
	for {
		// Read the line
		line, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		// build header map for later use.
		if lineNumber == 0 {
			for i, column := range line {
				headerMap[column] = i
			}
		} else if len(line) == len(headerMap) {

		} else {
			log.Fatalf(("CSV line %d does not have correct format."), lineNumber)
		}

		lineNumber++

		data = append(data, line)
	}

	return data, headerMap, nil
}

// getLocalIPs - returns a list of local IP addresses for the machine running the program
func getLocalIPs() ([]string, error) {
	var ips []string
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				ips = append(ips, ipNet.IP.String())
			}
		}
	}

	return ips, nil
}

// vensimTraffic - reads the CSV file and find ports that the local machine ip is the dst for and sends traffic where it is the src
func vensimTraffic(csvFile string) {

	ips, err := getLocalIPs()
	if err != nil {
		log.Fatal(err)
	}

	var ct clientTraffic
	// Open CSV File and create the reader and get data.
	data, header, err := LoadCSV(csvFile)
	if err != nil {
		log.Fatal(err)
	}

	uniqueConnections := make(map[string]bool)
	// Iterate through the different ips on the local machine
	for _, ip := range ips {
		// Iterate through the
		for _, row := range data {
			dst := row[header["dst_ip"]]
			port := row[header["port"]]
			protocol := row[header["proto"]]
			src := row[header["src_ip"]]

			//make a key to find duplicates using dst, port, and protocol
			key := fmt.Sprintf("%s:%s:%s", dst, port, protocol)

			if _, exists := uniqueConnections[key]; !exists {

				if ip == dst && protocol == "6" {
					ct.listenerTCPPorts = append(ct.listenerTCPPorts, port)
				} else if ip == dst && protocol == "17" {
					ct.listenerUDPPorts = append(ct.listenerUDPPorts, port)
				} else if ip == src && protocol == "6" {
					ct.dstTCPconnection = append(ct.dstTCPconnection, fmt.Sprintf("%s:%s", dst, port))
				} else if ip == src && protocol == "17" {
					ct.dstUDPconnection = append(ct.dstUDPconnection, fmt.Sprintf("%s:%s", dst, port))
				} else {
					continue
				}
				uniqueConnections[key] = true
			}
		}
	}

	//Tell user if the ips in the CSV dont match with the local machine
	if len(uniqueConnections) == 0 {
		log.Println("No listeners or connections found for this host")
		os.Exit(0)
	} else {
		log.Printf("Listeners: TCP: %v, UDP: %v\n", ct.listenerTCPPorts, ct.listenerUDPPorts)
		log.Printf("Connections: TCP: %v, UDP: %v\n", ct.dstTCPconnection, ct.dstUDPconnection)
	}

	// Create listeners and connections
	for _, port := range ct.listenerTCPPorts {
		go startTCPListener(port)
	}
	for _, port := range ct.listenerUDPPorts {
		go startUDPListener(port)
	}

	for _, addr := range ct.dstTCPconnection {
		go startTCPClient(addr)
	}
	for _, addr := range ct.dstUDPconnection {
		go startUDPClient(addr)
	}
	for {
		time.Sleep(1 * time.Second)
	}
}

// hostname - returns the hostname of the machine
func hostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "could_not_get_hostname"
	}
	return hostname
}

// openListeningPort - opens a listening port on the machine
func openListeningPort(port int) {
	flag.Parse()
	handler := http.HandlerFunc(handleRequest)
	http.Handle("/", handler)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)

}

// handleRequest - handles the request and sends a response back to the client
func handleRequest(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/text")
	w.Write([]byte(fmt.Sprintf("connected to %s from %s\r\n", hostname(), r.RemoteAddr)))
}

// ClearBOM - removes the BOM from the CSV file
func ClearBOM(r io.Reader) io.Reader {
	buf := bufio.NewReader(r)
	b, err := buf.Peek(3)
	if err != nil {
		// not enough bytes
		return buf
	}
	if b[0] == 0xef && b[1] == 0xbb && b[2] == 0xbf {
		buf.Discard(3)
	}
	return buf
}

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

// startUDPListener - starts up the server on the UDP port.  Sends the response to handleUDPPacket{
func startUDPListener(port string) {

	addr := fmt.Sprintf(":%s", port)
	conn, err := net.ListenPacket("udp", addr)
	if err != nil {
		log.Fatalf("Error starting UDP listener on port %s: %v", port, err)
	}
	defer conn.Close()

	log.Printf("UDP listener started on port %s", port)

	buffer := make([]byte, 1024)
	for {
		n, clientAddr, err := conn.ReadFrom(buffer)
		if err != nil {
			log.Printf("Error reading UDP packet: %v", err)
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
	log.Printf("UDP Listener responded with: %s to ip: %s", string(data), ip)
}

// randomInterval returns a random duration between min and max seconds
func randomInterval(min, max int) time.Duration {
	return time.Duration(min+rand.Intn(max-min+1)) * time.Second
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

// startUDPClient - starts the UDP client and sends a message to the server
func startUDPClient(addr string) {
	var conn net.Conn
	var err error

	//Open up the UDP port and clean up if something goes wrong
	for {
		conn, err = net.Dial("udp", addr)
		if err != nil {
			log.Printf("Error connecting to UDP server %v", err)
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
			log.Printf("Error reading from UDP server: %v", err)
		}
		responseTxt := string(buffer[:n])
		if len(string(buffer[:n])) == 0 {
			responseTxt = "No response"
		}
		log.Printf("UDP Connector response from %s: %s\n", addr, responseTxt)
		time.Sleep((randomInterval(MINTIME, MAXTIME)))
	}
}
