package diskplayer

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
)

func Record(url string) error {
	s, e := createSpotifyUri(url)
	if e != nil {
		return e
	}

	e = writeToDisk(s)
	if e != nil {
		return e
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
	p := GetConfigString(RECORD_PATH)

	b := []byte(spotifyUri)
	e := ioutil.WriteFile(p, b, 0644)
	if e != nil {
		return e
	}

	return nil
}
