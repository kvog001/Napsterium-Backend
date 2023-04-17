package dlservice

import (
	//"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var conn *websocket.Conn          // WebSocket connection
var urlChannel = make(chan string) // // Channel to receive youtube URLs from external source

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

	// Read messages from client
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Failed to read message:", err)
			break
		}
		log.Println("Received message:", string(message))
	}
}

func DownloadSong(youtubeURL string) {
	urlChannel <- youtubeURL
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
/*
var conn *websocket.Conn
var isWebSocketConnSetup bool
var urlChannel = make(chan string)
var connMutex sync.Mutex // Mutex to synchronize access to conn variable

func DownloadSong(youtubeURL string) {
    connMutex.Lock()
    defer connMutex.Unlock()

    if isWebSocketConnSetup {
        urlChannel <- youtubeURL
    } else {
        log.Println("WebSocket connection not set up yet. Discarding URL:", youtubeURL)
    }
}

func handleIncomingTitles() {
	for {
		url := <-urlChannel
		log.Println("Received youtube url from external source:", url)

		// Send the youtube url as a WebSocket message
		err := conn.WriteMessage(websocket.TextMessage, []byte(url))
		if err != nil {
			log.Fatal("Failed to send file title over WebSocket:", err)
		}
	}
}

func SetupWebsocketConn() {
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		// Upgrade the HTTP connection to a WebSocket connection
        	//connMutex.Lock()
        	conn, err := websocket.Upgrade(w, r, nil, 1024, 1024)
        	//connMutex.Unlock()
		
		log.Println("upgrading connection to websocket")
		
		if err != nil {
			log.Fatal("Failed to upgrade to WebSocket:", err)
		}
		defer conn.Close()

		// Set isWebSocketConnSetup to true after successful WebSocket upgrade
        	connMutex.Lock()
        	isWebSocketConnSetup = true
        	connMutex.Unlock()

		// Read the file content from the WebSocket response
		_, fileContent, err := conn.ReadMessage()
		if err != nil {
			log.Fatal("Failed to read file content:", err)
		}

		// Process the file content as needed
		fmt.Println("File content received from local server:", string(fileContent))
	})

	// Start the separate goroutine to handle incoming titles
	go handleIncomingTitles()

	log.Println("Listening to 8080")
	log.Fatal(http.ListenAndServe("193.233.202.119:8080", nil))
}
*/
