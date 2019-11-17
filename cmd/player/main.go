package main

import (
	"flag"
	"github.com/dinofizz/diskplayer"
	"log"
)


func main() {
	uri := flag.String("uri", "", "Spotify URI of album/playlist to play.")
	path := flag.String("path", "", "Path to file containing Spotify URI to play.")
	pause := flag.Bool("pause", false, "Pause Spotify playback.")
	flag.Parse()
	a := flag.Args()
	if len(a) != 0 {
		log.Fatalf("Unknown argument: %s. You might be missing a \"-\".", a[0]) // Expect user to eliminate unknown arguments
	}

	if *pause && (*uri != "" || *path != "") {
		flag.Usage()
		log.Fatal("Please specify either [pause] OR ONE OF [uri, path].")
	}

	if *uri != "" && *path != "" {
		flag.Usage()
		log.Fatal("Please specify either [uri] or [path], but not both.")
	}

	diskplayer.ReadConfig(diskplayer.DEFAULT_CONFIG_NAME)

	c, err := diskplayer.NewClient()
	if err != nil {
		log.Fatal(err)
	}

	if *pause {
		err = diskplayer.Pause(c)
	} else if *uri != "" {
		err = diskplayer.PlayUri(c, *uri)
	} else if *path != "" {
		err = diskplayer.PlayPath(c, *path)
	} else {
		err = diskplayer.Play(c)
	}

	if err != nil {
		log.Fatal(err)
	}
}
