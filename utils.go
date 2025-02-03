package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"math/rand"
	"net"
	"os"
	"strings"
	"time"
)

// hostIPAddresses - returns a list of local IP addresses for the machine running the program
func hostIPAddresses() ([]string, error) {
	var ips []string
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, fmt.Errorf("hostIPAddresses - %s", err)
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

// hostMatch takes an IP or an fqdn and returns true if it matches the host excecuting.
func hostMatch(host string) bool {
	ips, err := hostIPAddresses()
	if err != nil {
		LogErrorf("matchhost - %s", err)
	}

	// Check each ip
	for _, ip := range ips {
		if ip == host {
			return true
		}
	}
	// Check hostname
	return hostname() == host || strings.Split(hostname(), ".")[0] == host
}

// hostname - returns the hostname of the machine
func hostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "could_not_get_hostname"
	}
	return hostname
}

func LoadCSV(csvFile string) ([][]string, map[string]int, error) {
	// Open CSV File and create the reader
	file, err := os.Open(csvFile)
	if err != nil {
		LogErrorf("error opening csv file - %s", err)
	}
	defer file.Close()
	reader := csv.NewReader(ClearBOM(bufio.NewReader(file)))
	lineNumber := 0

	var data [][]string
	// Iterate through the CSV

	headerMap := make(map[string]int)
	for {
		// Increment line counter
		lineNumber++

		// Read the line
		line, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			LogErrorf("error reading csv file - %s", err)
		}

		// build header map for later use.
		if lineNumber == 1 {
			for i, header := range line {
				if header == "src_ip" {
					header = "src"
				}
				if header == "dst_ip" {
					header = "dst"
				}
				if header == "protocol" {
					header = "proto"
				}
				headerMap[header] = i
			}
		} else if len(line) != len(headerMap) {
			return nil, nil, fmt.Errorf("csv line %d - incorrect format. line has %d entries and headers have %d", lineNumber, len(line), len(headerMap))
		}

		data = append(data, line)
	}

	return data, headerMap, nil
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

// randomInterval returns a random duration between min and max seconds
func randomInterval(min, max int) time.Duration {
	return time.Duration(min+rand.Intn(max-min+1)) * time.Second
}
