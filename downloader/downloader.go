package downloader

import (
	"os"
	"fmt"
	"log"
	"bytes"
	"os/exec"
)

// e.g. youtubeURL := "https://www.youtube.com/watch?v=dQw4w9WgXcQ"

func DownloadSongToDisk(youtubeURL string, videoID string) {
	// Create the songsOpus directory if it doesn't already exist
	err := os.Mkdir("songsOpus", 0755)
	if err != nil && !os.IsExist(err) {
		fmt.Println("Error creating songs directory.")
		return
	}

	// Download the opus file using yt-dlp
	cmd := exec.Command(
	"yt-dlp", 
	"--extract-audio", 
	"--audio-format", "opus", 
	"--audio-quality", "192k", 
	"-o", "songsOpus/" + videoID + ".%(ext)s", youtubeURL)

	var stdout, stderr bytes.Buffer
        cmd.Stdout = &stdout
        cmd.Stderr = &stderr

	// Run the command and wait for it to finish
	err = cmd.Run()
	if err != nil {
		log.Printf("Error downloading opus file: %s", err)
		log.Printf("yt-dlp command output: %s", stdout.String())
    		log.Printf("yt-dlp command error message: %s", stderr.String())
		return
	}
	log.Printf("File downloaded successfully at songsOpus/%s\n", videoID)

	logFileSize(videoID)
}

func logFileSize(videoID string) {
	// Get the file info for the downloaded opus file
	filePath := "songsOpus/" + videoID + ".opus"
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		log.Printf("Error getting file info: %s", err)
		return
	}

	// Get the size of the downloaded opus file in bytes
	fileSize := fileInfo.Size()

	// Convert the file size to kilobytes
	fileSizeKB := float64(fileSize) / (1024)

	log.Printf("Downloaded opus file size: %.2f KB\n", fileSizeKB)
}
