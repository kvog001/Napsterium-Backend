package main

import (
	"os"
	"log"
	"fmt"
	"net/http"
  "io/ioutil"
  "github.com/gorilla/websocket"
	"Napsterium-Backend/downloader"
)

func main() {
	// Define the WebSocket dialer with appropriate configuration
	var dialer = websocket.Dialer{
		Proxy:            http.ProxyFromEnvironment,
		HandshakeTimeout: 5000, // Set an appropriate handshake timeout
	}

	// Create the WebSocket connection to your Golang web server
	conn, _, err := dialer.Dial("ws://kvogli.xyz:8080/ws", nil) // Update with your Golang web server URL
	if err != nil {
		log.Fatal("Failed to connect to web server:", err)
	}
	defer conn.Close()

	// Read the WebSocket message (YouTube URL)
	_, data, err := conn.ReadMessage()
	if err != nil {
		log.Fatal("Failed to read WebSocket message:", err)
	}

	youtubeURL := string(data)

	// download song
	downloader.DownloadSongToDisk(youtubeURL, "mp3", "7")

	// Use the song id to read the file on the local server
	songID := downloader.ExtractSongID(youtubeURL)
	file, err := os.Open(string(songID + ".mp3"))
	if err != nil {
		log.Fatal("Failed to open file:", err)
	}
	defer file.Close()

	// Read the file content and send it as a response back to the web server
	fileContent, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal("Failed to read file:", err)
	}

	err = conn.WriteMessage(websocket.TextMessage, fileContent)
	if err != nil {
		log.Fatal("Failed to send file content:", err)
	}

	fmt.Println("File content sent to web server.")
}

