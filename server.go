package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"crypto/tls"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/helloworld", helloHandler)

	fmt.Printf("Starting server at port 8080\n")
	// Create the TLS config
    	config := &tls.Config{
        	MinVersion: tls.VersionTLS10,
        	MaxVersion: tls.VersionTLS13,
    	}

    	// Create the HTTP server with the TLS config
    	server := &http.Server{
       		Addr:      "193.233.202.119:8080",
		Handler:  mux,
        	TLSConfig: config,
    	}

    	// Listen and serve with TLS
    	err := server.ListenAndServeTLS("server.crt", "server.key")
    	if err != nil {
		fmt.Println("Error starting server:", err)
    	}

	//if err := http.ListenAndServe("193.233.202.119:8080", nil); err != nil {
	//	log.Fatal(err)
	//}
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

