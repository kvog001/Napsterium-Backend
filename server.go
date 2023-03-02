package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/helloworld", helloHandler)

	fmt.Printf("Starting server at port 443\n")

	cert, err := tls.LoadX509KeyPair("/etc/letsencrypt/live/kvogli.xyz/fullchain.pem", "/etc/letsencrypt/live/kvogli.xyz/privkey.pem")
	if err != nil {
		log.Fatalf("Failed to load SSL certificate: %v", err)
	}
	// Create the TLS config
	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS10,
	}

	addr := "0.0.0.0:443"//"193.233.202.119:443"
	// Create the HTTP server with the TLS config
	server := &http.Server{
		Addr:      addr,
		Handler:   mux,
		TLSConfig: config,
	}

	// Listen and serve with TLS
	log.Printf("Listening on %s...\n", addr)
	err = server.ListenAndServeTLS("", "")
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/helloworld" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	// Read the request body into a string
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body.", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()
	requestText := string(body)

	// Print the request text to the server console
	fmt.Printf("Received request: %s\n", requestText)

	// Return a response
	fmt.Fprintf(w, "Hello World!\n")
}
