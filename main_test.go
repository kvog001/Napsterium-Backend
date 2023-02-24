package main

import (
	"os"
	"testing"
)

func TestDownloadSongMP3(t *testing.T) {
	// call the function to download a song
	youtubeLink := "https://www.youtube.com/watch?v=dQw4w9WgXcQ"
	downloadSongMP3(youtubeLink)

	// check if the file was downloaded
	_, err := os.Stat("/songs/Rick Astley - Never Gonna Give You Up (Video).mp3")
	if os.IsNotExist(err) {
		t.Errorf("Downloaded file does not exist: %v", err)
	}
}
