package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"os"
	"time"
)

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
		log.Fatal(err)
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
			log.Fatal(err)
		}

		// build header map for later use.
		if lineNumber == 0 {
			for i, column := range line {
				if column == "src_ip" {
					column = "src"
				}
				if column == "dst_ip" {
					column = "dst"
				}
				headerMap[column] = i
			}
		} else if len(line) != len(headerMap) {
			return nil, nil, fmt.Errorf("csv line %d - incorrect format", lineNumber)
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
