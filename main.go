package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

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

	// Call open
	if os.Args[1] == "open" {
		port, err := strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatalf("cannot process port - %s", err)

		}
		openListeningPort(port)
	}

	// Call traffic
	if os.Args[1] == "traffic" {
		// Validate hostname
		fmt.Printf("hostname discovered as %s\r\n", hostname())
		sendTraffic(os.Args[2])
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

func hostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "could_not_get_hostname"
	}
	return hostname
}

func openListeningPort(port int) {
	flag.Parse()
	handler := http.HandlerFunc(handleRequest)
	http.Handle("/", handler)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)

}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/text")
	w.Write([]byte(fmt.Sprintf("connected to %s from %s\r\n", hostname(), r.RemoteAddr)))
}
