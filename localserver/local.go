package main

import (
	"log"
	"fmt"
  	"github.com/gorilla/websocket"
)

func main() {
	// Dial WebSocket server
	conn, _, err := websocket.DefaultDialer.Dial("ws://kvogli.xyz:8080/ws", nil)
	if err != nil {
		log.Fatal("Failed to connect to WebSocket server:", err)
	}
	defer conn.Close()

	// Send message to server
	message := "Hello, server!"
	err = conn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		log.Fatal("Failed to send message to server:", err)
	}

	// Start a separate goroutine to continuously read messages from the server
	go func() {
		for {
			_, response, err := conn.ReadMessage()
			if err != nil {
				log.Fatal("Failed to read response from server:", err)
			}
			fmt.Println("Received response from server:", string(response))
		}
	}()

	// Wait for user input to exit
	fmt.Println("Press enter to exit...")
	fmt.Scanln()
}

