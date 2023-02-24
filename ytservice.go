package main

import (
	"fmt"
	"os/exec"
)

// e.g. youtubeLink := "https://www.youtube.com/watch?v=dQw4w9WgXcQ"
func downloadSongMP3(youtubeLink string) {
	fmt.Println("downloading a song")
	cmd := exec.Command("youtube-dl", "--verbose","--extract-audio", "--audio-format", "mp3", "-o", "/songs/%(title)s.%(ext)s", youtubeLink)
	out, err := cmd.CombinedOutput()
	if err != nil {
			fmt.Println("Error:", err)
	}
	fmt.Println(string(out))
}