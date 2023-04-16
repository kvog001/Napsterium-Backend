package main

import (
	"fmt"
	"log"
	"time"
	"net/http"
	"crypto/tls"
	"Napsterium-Backend/handler"
	"Napsterium-Backend/dbservice"
	"Napsterium-Backend/dlservice"
)

func main() {
	dbservice.ConnectToDB()

	go dlservice.SetupWebsocketConn()

	// Start the download after 5 seconds
	go func() {
		time.Sleep(5 * time.Second)
		dlservice.DownloadSong("https://www.youtube.com/watch?v=03qttyOlvQE")
	}()

	mux := http.NewServeMux()
	mux.HandleFunc("/helloworld", handler.HelloHandler)

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

	addr :=  "193.233.202.119:443" /* "0.0.0.0:443" */
	// Create the HTTP server with the TLS config
	server := &http.Server {
		Addr:      addr,
		Handler:   mux,
		TLSConfig: config,
	}

	// Listen and serve with TLS
	log.Printf("Listening on %s...\n", addr)
	err = server.ListenAndServeTLS("", "")
	if err != nil {
		log.Println("Error starting server:", err)
	}

}
