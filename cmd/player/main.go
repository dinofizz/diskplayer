package main

import (
	"flag"
	"github.com/dinofizz/diskplayer/internal/config"
	"github.com/dinofizz/diskplayer/internal/spotifyplayer"
	"log"
)

func main() {
	uri := flag.String("uri", "", "Spotify URI of album/playlist to play.")
	pause := flag.Bool("pause", false, "Pause Spotify playback.")
	flag.Parse()

	if *pause && *uri != "" {
		log.Fatal("Please specify either [pause] OR [uri], but not both.")
	} else if !*pause && *uri == "" {
		log.Fatal("Spotify URI is required for playback.")
	}

	config.ReadConfig()

	if *pause {
		spotifyplayer.Pause()
	} else {
		spotifyplayer.Play(*uri)
	}
}
