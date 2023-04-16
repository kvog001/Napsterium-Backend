package downloader

import (
	"os"
	"fmt"
	"log"
	"bytes"
	"net/url"
	"os/exec"
)

const SongsPath = "songs"

func DownloadSongToDisk(youtubeURL, audioFormat, audioQuality string) {
	// Create the songs directory if it doesn't already exist
	err := os.Mkdir(SongsPath, 0755)
	if err != nil && !os.IsExist(err) {
		fmt.Println("Error creating songs directory.")
		return
	}

	songID := ExtractSongID(youtubeURL)
	executeYTDLP(songID, youtubeURL, audioFormat, audioQuality)
	logFileSize(songID)
}

func executeYTDLP(songID, youtubeURL, audioFormat, audioQuality string) {
	// Download the song file using yt-dlp
	cmd := exec.Command(
		"yt-dlp", 
		"--extract-audio", 
		"--audio-format", audioFormat, 
		"--audio-quality", audioQuality, 
		"-o", SongsPath + "/" + songID + ".%(ext)s", youtubeURL)
	
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Run the command and wait for it to finish
	err := cmd.Run()
	if err != nil {
		log.Printf("yt-dlp error downloading file: %s", err)
		log.Printf("yt-dlp command output: %s", stdout.String())
		log.Printf("yt-dlp command error message: %s", stderr.String())
		return
	}
	log.Printf("File downloaded successfully at %s/%s\n", SongsPath, songID)
}

func ExtractSongID(youtubeURL string) string {
	u, err := url.Parse(youtubeURL)
	if err != nil {
		return ""
	}

	query := u.Query()
	return query.Get("v")
}

func logFileSize(songID string) {
	// Get the file info for the downloaded file
	filePath := SongsPath + "/" + songID + ".mp3"
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		log.Printf("Error getting file info: %s", err)
		return
	}

	// Get the size of the downloaded file in bytes
	fileSize := fileInfo.Size()

	// Convert the file size to kilobytes
	fileSizeKB := float64(fileSize) / (1024)

	log.Printf("Downloaded file size: %.2f KB\n", fileSizeKB)
}
