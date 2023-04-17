package main

import (
	"log"
	"fmt"
	"io/ioutil"
  "github.com/gorilla/websocket"
	"Napsterium-Backend/downloader"
)

func main() {
	// Dial WebSocket server
	conn, _, err := websocket.DefaultDialer.Dial("ws://kvogli.xyz:8080/ws", nil)
	if err != nil {
		log.Fatal("Failed to connect to WebSocket server:", err)
	}
	defer conn.Close()

	// Continuously read messages from the server
	for {
		_, request, err := conn.ReadMessage()
		if err != nil {
			log.Fatal("Failed to read request from server:", err)
		}

		youtubeURL := string(request)
		fmt.Println("Received request from server:", youtubeURL)
		
		songID := downloader.ExtractSongID(youtubeURL)
		downloader.DownloadSongToDisk(youtubeURL)

		data, err := ioutil.ReadFile(downloader.SongsPath + "/" + songID + "." + downloader.DownloadFormat)
		if err != nil {
			log.Println("Error reading/loading song file from disk.")
			return
		}

		// Send response to server
		err = conn.WriteMessage(websocket.TextMessage, []byte(data))
		if err != nil {
			log.Println("Failed to send response to server:", err)
		}
		// TODO: delete song after sending it
	}
}
