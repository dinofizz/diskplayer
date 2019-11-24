package main

import (
	"flag"
	"github.com/dinofizz/diskplayer"
	"golang.org/x/oauth2"
	"log"
	"os"
)

func main() {
	auth := flag.Bool("auth", false, "Retrieve a new Spotify OAuth2 token.")
	uri := flag.String("uri", "", "Spotify URI of album/playlist to play.")
	path := flag.String("path", "", "Path to file containing Spotify URI to play.")
	pause := flag.Bool("pause", false, "Pause Spotify playback.")
	flag.Parse()
	a := flag.Args()
	if len(a) != 0 {
		log.Fatalf("Unknown argument: %s. You might be missing a \"-\".", a[0]) // Expect user to eliminate unknown arguments
	}

	if (*auth && *pause) || (*auth && (*uri != "" || *path != "")) || (*pause && (*uri != "" || *path != "")) {
		flag.Usage()
		log.Fatal("Please specify either [auth] OR [pause] OR ONE OF [uri, path].")
	}

	if *uri != "" && *path != "" {
		flag.Usage()
		log.Fatal("Please specify either [uri] or [path], but not both.")
	}

	diskplayer.ReadConfig(diskplayer.DEFAULT_CONFIG_NAME)

	an, err := diskplayer.NewAuthenticator()
	if err != nil {
		log.Fatal(err)
	}

	if *auth {
		ch := make(chan *oauth2.Token, 1)
		s := diskplayer.NewDiskplayerServer(an, ch)

		t, err := diskplayer.NewToken(s)
		if err != nil {
			log.Fatal(err)
		}
		err = diskplayer.SaveToken(t)
		if err != nil {
			log.Fatal(err)
		}

		os.Exit(0)
	}

	t, err := diskplayer.ReadToken()
	if err != nil {
		log.Fatal(err)
	}

	c := diskplayer.NewClient(an, t)

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
