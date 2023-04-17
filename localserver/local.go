package main

import (
	"log"
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
		log.Println("Received request from server:", youtubeURL)
		
		songID := downloader.ExtractSongID(youtubeURL)
		downloader.DownloadSongToDisk(youtubeURL)

		songPath := downloader.SongsPath + "/" + songID + "." + downloader.DownloadFormat
		log.Printf("reading song data from disk %s ...\n", songPath)
		data, err := ioutil.ReadFile(songPath)
		if err != nil {
			log.Println("Error reading/loading song file from disk.")
			return
		}

		log.Println("sending response to server ...")
		// Send response to server
		err = conn.WriteMessage(websocket.TextMessage, []byte(data))
		if err != nil {
			log.Println("Failed to send response to server:", err)
		}
		// TODO: delete song after sending it
	}
}
