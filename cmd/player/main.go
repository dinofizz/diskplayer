package main

import (
	"flag"
	"github.com/dinofizz/diskplayer/internal/spotifyplayer"
	"log"
	"os"
)

func main() {
	deviceName := flag.String("device", "", "Name of Spotify player device.")
	playUri := flag.String("uri", "", "Spotify URI of album/playlist to play.")
	pause := flag.Bool("pause", false, "Pause Spotify playback.")
	flag.Parse()

	if *pause && (*deviceName != "" || *playUri != "") {
		log.Fatal("Please specify either [pause] OR [device and uri], but not combinations of both sets of flags.")
	}

	if *pause {
		spotifyplayer.Pause()
		os.Exit(0)
	}

	if *deviceName != "" && *playUri != "" {
		spotifyplayer.Play(*deviceName, *playUri)
	} else {
		log.Fatal("Both device name and Spotify URI are required for playback.")
	}
}
