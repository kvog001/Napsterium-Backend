package dlservice

import (
	//"fmt"
	"log"
	"net/url"
	"net/http"
	"github.com/gorilla/websocket"
)

var conn *websocket.Conn            // WebSocket connection
var urlChannel = make(chan string) // Channel to receive youtube URLs from external source
var songChannel = make(chan []byte)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow any origin to connect
	},
}

func WebsocketHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade HTTP connection to WebSocket
	var err error
	conn, err = upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("Failed to upgrade to WebSocket:", err)
	}
	defer conn.Close()

	go handleOutgoingMessages() // Start a new goroutine to handle outgoing messages

	// Read song audio data from local
	for {
		_, song, err := conn.ReadMessage()
		if err != nil {
			log.Println("Failed to read message:", err)
			break
		}
		
		songChannel <- song
		//log.Println("Received message:", string(message))
	}
}

func DownloadSong(youtubeURL string) []byte {
	urlChannel <- youtubeURL
	songWebm := <-songChannel
	songMP3 := convert_webm_to_mp3(songWebm)
	return songMP3
}

func handleOutgoingMessages() {
	for {
		url := <-urlChannel
		log.Println("Received youtube url from external source:", url)
		// Send the youtube url as a WebSocket message
		err := conn.WriteMessage(websocket.TextMessage, []byte(url))
		if err != nil {
			log.Println("Failed to send file title over WebSocket:", err)
		}
	}
}

func ExtractSongID(youtubeURL string) string {
	u, err := url.Parse(youtubeURL)
	if err != nil {
		return ""
	}

	query := u.Query()
	return query.Get("v")
}
