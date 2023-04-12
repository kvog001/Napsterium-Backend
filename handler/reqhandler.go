package handler

import (
	"fmt"
	"log"
	"net/url"
	"net/http"
	"io/ioutil"
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

	downloadSongToPath(youtubeURL, videoID)

	// Return a response
	log.Println("Preparing response")
	// Load the mp3 file from disk
	fileBytes, err := ioutil.ReadFile("songsMP3/" + videoID + ".mp3")
	if err != nil {
		http.Error(w, "Error reading mp3 file.", http.StatusInternalServerError)
		return
	}

	// Set the content type header to indicate that we're returning binary data
	w.Header().Set("Content-Type", "application/octet-stream")

	// Set the content disposition header to suggest a filename
	w.Header().Set("Content-Disposition", "attachment; filename=" + videoID + ".mp3")

	// Write the mp3 file bytes to the response writer
	w.Write(fileBytes)
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
