package diskplayer

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
)

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

func writeToDisk(spotifyUri string) error {
	p := ConfigValue(RECORD_PATH)

	b := []byte(spotifyUri)
	err := ioutil.WriteFile(p, b, 0644)
	if err != nil {
		return err
	}

	return nil
}
