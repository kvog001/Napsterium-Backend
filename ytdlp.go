package main

import (
	"os"
	"fmt"
	"log"
	"bytes"
	"os/exec"
)

// e.g. youtubeURL := "https://www.youtube.com/watch?v=dQw4w9WgXcQ"

func downloadSongToPath(youtubeURL string, videoID string) {
	// Create the songsMP3 directory if it doesn't already exist
	err := os.Mkdir("songsMP3", 0755)
	if err != nil && !os.IsExist(err) {
		fmt.Println("Error creating songs directory.")
		return
	}

	// Download the mp3 file using yt-dlp
	cmd := exec.Command("yt-dlp", "--extract-audio", "--audio-format", "mp3", "-o", "songsMP3/" + videoID + ".%(ext)s", youtubeURL)

	var stdout, stderr bytes.Buffer
        cmd.Stdout = &stdout
        cmd.Stderr = &stderr

	// Run the command and wait for it to finish
	err = cmd.Run()
	if err != nil {
		log.Printf("Error downloading mp3 file: %s", err)
		log.Printf("yt-dlp command output: %s", stdout.String())
    		log.Printf("yt-dlp command error message: %s", stderr.String())
		return
	}
	log.Printf("File downloaded successfully at songsMP3/%s\n", videoID)
}