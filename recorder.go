package diskplayer

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
)

// Record takes in a web URL which links to a Spotify album or playlist and records the corresponding Spotify ID to
// the filepath specified in the diskplayer.yaml configuration file under the recorder.file_path field.
// The web URL should be something like https://open.spotify.com/album/1S7mumn7D4riEX2gVWYgPO
// Returns an error if one is encountered.
func Record(url string) error {
	s, err := createSpotifyUri(url)
	if err != nil {
		return err
	}

	err = writeToDisk(s)
	if err != nil {
		return err
	}

	return nil
}

// createSpotifyUri creates a Spotify URI from the web URL.
// The web URL should be something like https://open.spotify.com/album/1S7mumn7D4riEX2gVWYgPO
// A string representing the Spotify URI is returned or any error that is encountered.
func createSpotifyUri(url string) (string, error) {
	i := strings.LastIndex(url, "/")
	id := url[i+1:]
	var u string
	if strings.Contains(url, "/album/") {
		u = "spotify:album:" + id
	} else if strings.Contains(url, "/playlist/") {
		u = "spotify:playlist:" + id
	} else {
		return "", errors.New(fmt.Sprintf("URL represents neither album nor playlist: %s", url))
	}
	return u, nil
}

// writeToDisk takes a string containing a Spotify URI and writes to the the filepath specified in the diskplayer.yaml
// configuration file under the recorder.file_path field.
// Returns an error if one is encountered.
func writeToDisk(spotifyUri string) error {
	p := ConfigValue(RECORD_PATH)

	b := []byte(spotifyUri)
	err := ioutil.WriteFile(p, b, 0644)
	if err != nil {
		return err
	}

	return nil
}
