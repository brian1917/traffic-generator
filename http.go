package main

import (
	"flag"
	"fmt"
	"net/http"
)

// openHttpListener - opens a listening port on the machine
func openHttpListener(port int) {
	flag.Parse()
	handler := http.HandlerFunc(handleHttpRequest)
	http.Handle("/", handler)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)

}

// handleHttpRequest - handles the request and sends a response back to the client
func handleHttpRequest(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/text")
	w.Write([]byte(fmt.Sprintf("connected to %s from %s\r\n", hostname(), r.RemoteAddr)))
}
