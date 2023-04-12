package handler

import (
	"fmt"
	"log"
	"net/url"
	"net/http"
	"io/ioutil"
	"Napsterium-Backend/downloader"
	"github.com/hraban/opus"
)

func HelloHandler(w http.ResponseWriter, r *http.Request) {
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
	youtubeURL := string(body)

	// Print the requested YouTube URL to the server console
	fmt.Printf("Received request: %s\n", youtubeURL)

	videoID := extractVideoID(youtubeURL)

	downloader.DownloadSongToDisk(youtubeURL, videoID)

	// Return a response
	log.Println("Preparing response")
	// Load the opus file from disk
	fileBytes, err := ioutil.ReadFile("songsOpus/" + videoID + ".opus")
	if err != nil {
		http.Error(w, "Error reading opus file.", http.StatusInternalServerError)
		return
	}

	// decode .opus to pcm
	const channels = 1
	const sampleRate = 48000

	dec, err := opus.NewDecoder(sampleRate, channels)
	if err != nil {
			fmt.Println("Failed to create Opus decoder:", err)
			return
	}
	var frameSizeMs = float32(60)  // if you don't know, go with 60 ms.
	frameSize := channels * frameSizeMs * sampleRate / 1000
	pcm := make([]int16, int(frameSize))
	n, err := dec.Decode(fileBytes, pcm)
	if err != nil {
			fmt.Println("Failed to decode Opus file:", err)
			return
	}

	pcm = pcm[:n * channels]

	// Set the content type header to indicate that we're returning binary data
	w.Header().Set("Content-Type", "application/octet-stream")

	// Set the content disposition header to suggest a filename
	w.Header().Set("Content-Disposition", "attachment; filename=" + videoID + ".opus")

	// Write the mp3 file bytes to the response writer
	w.Write(int16ToBytes(pcm))
	log.Println("Returning response")
}

func extractVideoID(youtubeURL string) string {
	u, err := url.Parse(youtubeURL)
	if err != nil {
		return ""
	}

	query := u.Query()
	return query.Get("v")
}

// Utility function to convert []int16 to []byte
func int16ToBytes(data []int16) []byte {
	bytes := make([]byte, len(data)*2)
	for i, d := range data {
		bytes[i*2] = byte(d)
		bytes[i*2+1] = byte(d >> 8)
	}
	return bytes
}
