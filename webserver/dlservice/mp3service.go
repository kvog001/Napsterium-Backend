package dlservice

import (
	"bytes"
	"log"

	"github.com/hajimehoshi/go-mp3"
)

func convert_webm_to_mp3(webmData []byte) []byte {
	// Create a bytes.Reader from the .webm file data
	webmReader := bytes.NewReader(webmData)

	// Decode the .webm audio data
	decoder, err := mp3.NewDecoder(webmReader)
	if err != nil {
		log.Fatal(err)
	}

	// Create a bytes.Buffer to hold the .mp3 audio data
	mp3Buffer := new(bytes.Buffer)

	// Decode and write the .mp3 audio to the buffer
	buf := make([]byte, 8192)
	for {
		n, err := decoder.Read(buf)
		if err != nil {
			log.Fatal(err)
		}
		if n == 0 {
			break
		}
		_, err = mp3Buffer.Write(buf[:n])
		if err != nil {
			log.Fatal(err)
		}
	}

	// Get the .mp3 audio data as []byte from the buffer
	mp3Data := mp3Buffer.Bytes()
	return mp3Data
}