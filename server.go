package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
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

	addr := /* "193.233.202.119:443" */ "0.0.0.0:443"
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
	ytURL := string(body)

	// Print the requested YouTube URL to the server console
	fmt.Printf("Received request: %s\n", ytURL)

	videoID := extractVideoID(ytURL)

	downloadSongToPath(ytURL, videoID)

	// Return a response
	// Load the mp3 file from disk
	fileBytes, err := ioutil.ReadFile("songsMP3/" + videoID + ".mp3")
	if err != nil {
		http.Error(w, "Error reading mp3 file.", http.StatusInternalServerError)
		return
	}

	// Set the content type header to indicate that we're returning binary data
	w.Header().Set("Content-Type", "application/octet-stream")

	// Set the content disposition header to suggest a filename
	w.Header().Set("Content-Disposition", "attachment; filename="+videoID+".mp3")

	// Write the mp3 file bytes to the response writer
	w.Write(fileBytes)
}

func downloadSongToPath(ytURL string, videoID string) {
	// Create the songsMP3 directory if it doesn't already exist
	err := os.Mkdir("songsMP3", 0755)
	if err != nil && !os.IsExist(err) {
		fmt.Println("Error creating songs directory.")
		return
	}

	// Download the mp3 file using yt-dlp
	cmd := exec.Command("yt-dlp", "--extract-audio", "--audio-format", "mp3", "-o", "songsMP3/"+videoID+".%(ext)s", ytURL)
	// Run the command and wait for it to finish
	err = cmd.Run()
	if err != nil {
		log.Printf("Error downloading mp3 file: %s", err)
		return
	}
	log.Printf("File downloaded successfully at songsMP3/%s\n", videoID)
}

func extractVideoID(ytURL string) string {
	u, err := url.Parse(ytURL)
	if err != nil {
		return ""
	}

	query := u.Query()
	return query.Get("v")
}
