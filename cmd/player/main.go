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
	a := flag.Args()
	if len(a) != 0 {
		log.Fatalf("Unknown argument: %s. You might be missing a \"-\".", a[0]) // Expect user to eliminate unknown arguments
	}

	if *pause && *uri != "" {
		flag.Usage()
		log.Fatal("Please specify either [pause] OR [uri], but not both.")
	}

	config.ReadConfig()

	if *pause {
		spotifyplayer.Pause()
	} else if *uri != "" {
		spotifyplayer.PlayUri(*uri)
	} else {
		spotifyplayer.Play()
	}
}
