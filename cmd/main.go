package main

import (
	"log"
	"net/http"

	spotify "github.com/noahschumacher/spotify/internal"
)

func main() {

	// Start the http server
	http.HandleFunc("/authurl", spotify.AuthURLBuilder)
	http.HandleFunc("/callback", spotify.AuthCallback)
	http.HandleFunc("/usertrack", spotify.UserTrack)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Error:%s", err)
	}
}
