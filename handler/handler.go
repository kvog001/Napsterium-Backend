package handler

import (
	"fmt"
	"log"
	"net/http"
	"io/ioutil"
	"github.com/hraban/opus"
	"Napsterium-Backend/downloader"
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

	fmt.Printf("Received request: %s\n", youtubeURL)

	downloader.DownloadSongToDisk(youtubeURL, "mp3", "9")
	songID := downloader.ExtractSongID(youtubeURL)
	sendResponse(w, songID)
}

func sendResponse(w http.ResponseWriter, songID string) {
	log.Println("--- preparing Response ---")
	// Load the song file from disk
	data, err := ioutil.ReadFile(downloader.SongsPath + "/" + songID + ".mp3")
	if err != nil {
		http.Error(w, "Error reading/loading song file from disk.", http.StatusInternalServerError)
		return
	}

	// Set the content type header to indicate that we're returning binary data
	w.Header().Set("Content-Type", "application/octet-stream")

	// Set the content disposition header to suggest a filename
	w.Header().Set("Content-Disposition", "attachment; filename=" + songID + ".mp3")

	// Write the mp3 file bytes to the response writer
	w.Write(data)
	log.Println("--- sending Response ---")
}

func decodeOpusToPcm(data []byte) []byte {
	const channels = 1
	const sampleRate = 48000

	dec, err := opus.NewDecoder(sampleRate, channels)
	if err != nil {
			fmt.Println("Failed to create Opus decoder:", err)
			return []byte{}
	}
	var frameSizeMs = float32(60)  // if you don't know, go with 60 ms.
	frameSize := channels * frameSizeMs * sampleRate / 1000
	pcm := make([]int16, int(frameSize))
	n, err := dec.Decode(data, pcm)
	if err != nil {
			fmt.Println("Failed to decode Opus file:", err)
			return []byte{}
	}

	pcm = pcm[:n * channels]
	return int16ToBytes(pcm)
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
