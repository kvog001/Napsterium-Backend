package handler

import (
	"log"
	"net/http"
	"io/ioutil"
	"github.com/hraban/opus"
	"Napsterium-Backend/dlservice"
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

	log.Printf("Received request: %s\n", youtubeURL)

	songID := dlservice.ExtractSongID(youtubeURL)
	song := dlservice.DownloadSong(youtubeURL)

	sendResponse(w, songID, song)
}

func sendResponse(w http.ResponseWriter, songID string, data []byte) {
	log.Println("--- preparing Response ---")

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
		log.Println("Failed to create Opus decoder:", err)
		return []byte{}
	}
	var frameSizeMs = float32(60)  // if you don't know, go with 60 ms.
	frameSize := channels * frameSizeMs * sampleRate / 1000
	pcm := make([]int16, int(frameSize))
	n, err := dec.Decode(data, pcm)
	if err != nil {
		log.Println("Failed to decode Opus file:", err)
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
