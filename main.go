package main

import (
	"fmt"
	"io"
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

	// Make sure we have a command
	if len(os.Args) < 2 {
		showHelp()

		// Process open-listener command
	} else if os.Args[1] == "open-listener" {
		if len(os.Args) != 3 {
			showHelp()
			LogError("open-listener requires another argument for port (e.g., traffic-generator open 3306)")
		}
		port, err := strconv.Atoi(os.Args[2])
		if err != nil {
			LogErrorf("error converting port to int - %s", err)
		}
		openHttpListener(port)

		// Process send-traffic command
	} else if os.Args[1] == "send-traffic" {
		if len(os.Args) != 3 {
			showHelp()
			LogError("send-traffic requires another argument for the traffic import file.")
		}
		LogInfof(true, "hostname discovered as %s\r\n", hostname())
		sendTraffic(os.Args[2])

		// Process continuous command
	} else if os.Args[1] == "continuous" {
		if len(os.Args) != 3 {
			showHelp()
			LogError("continuous requires another argument for the import file.")
		}
		LogInfof(false, "hostname discovered as %s\r\n", hostname())
		openAndContinuousTraffic(os.Args[2])

		// Everything else is invalid
	} else {
		showHelp()
	}
}

func sendTraffic(csvFile string) {
	csvData, headers, err := LoadCSV(csvFile)
	if err != nil {
		LogErrorf("error loading csv file - %s", err)
	}

	// Iterate through the CSV
	for index, line := range csvData {

		// Skip header
		if index == 0 {
			continue
		}

		// Skip if it's not for this current host
		if line[headers["src"]] != hostname() {
			continue
		}

		// Make HTTP request for
		resp, err := http.Get(fmt.Sprintf("http://%s:%s", line[1], line[2]))
		if err != nil {
			LogErrorf("error making http request - %s", err)
		}
		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			LogErrorf("error reading response body - %s", err)
		}
		LogInfo(string(responseBody), true)
	}
}

// openAndContinuousTraffic - reads the CSV file and find ports that the local machine ip is the dst for and sends traffic where it is the src
func openAndContinuousTraffic(csvFile string) {

	ips, err := getLocalIPs()
	if err != nil {
		LogErrorf("error getting local ips - %s", err)
	}

	var ct clientTraffic
	// Open CSV File and create the reader and get data.
	data, header, err := LoadCSV(csvFile)
	if err != nil {
		LogErrorf("error loading csv file - %s", err)
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
		LogInfof(true, "no listeners or connections found for this host")
		os.Exit(0)
	} else {
		LogInfof(true, "listeners: tcp: %v, udp: %v\n", ct.listenerTCPPorts, ct.listenerUDPPorts)
		LogInfof(true, "connections: tcp: %v, udp: %v\n", ct.dstTCPconnection, ct.dstUDPconnection)
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
