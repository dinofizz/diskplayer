// This example demonstrates how to authenticate with Spotify using the authorization code flow.
// In order to run this example yourself, you'll need to:
//
//  1. Register an application at: https://developer.spotify.com/my-applications/
//       - Use "http://localhost:8080/callback" as the redirect URI
//  2. Set the SPOTIFY_ID environment variable to the client ID you got in step 1.
//  3. Set the SPOTIFY_SECRET environment variable to the client secret from step 1.
package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// redirectURI is the OAuth redirect URI for the application.
// You must register an application at Spotify's developer portal
// and enter this value.
const redirectURI = "http://localhost:8080/callback"

var (
	auth = spotify.NewAuthenticator(redirectURI, spotify.ScopeUserReadPrivate, spotify.ScopePlaylistReadPrivate,
		spotify.ScopeUserModifyPlaybackState, spotify.ScopeUserReadPlaybackState)
	ch    = make(chan *spotify.Client, 1)
	state = "abc123"
)

func main() {
	deviceName := flag.String("device", "", "Name of Spotify player device.")
	readPath := flag.String("path", "", "Path to file containing Spotify URI.")
	playUri := flag.String("uri", "", "Spotify URI of album/playlist to play.")
	pause := flag.Bool("pause", false, "Pause Spotify playback.")
	recordUrl := flag.String("record", "", "URL of item to be recorded.")
	flag.Parse()

	// first start an HTTP server
	http.HandleFunc("/callback", completeAuth)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got request for:", r.URL.String())
	})

	var server *http.Server

	token, err := tokenFromFile("tokenFile")
	if err != nil {
		server = &http.Server{Addr: ":8080", Handler: nil}
		go server.ListenAndServe()
		url := auth.AuthURL(state)
		fmt.Println("Please log in to Spotify by visiting the following page in your browser:", url)
	} else {
		client := auth.NewClient(token)
		ch <- &client
	}

	client := <-ch

	if server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := server.Shutdown(ctx)
		handleError(err)
	}

	if *pause {
		err := client.Pause()
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}

	if *recordUrl != "" {
		slashIndex := strings.LastIndex(*recordUrl, "/")
		id := (*recordUrl)[slashIndex+1:]

		var uri spotify.URI
		var uriString string

		if strings.Contains(*recordUrl, "/album/") {
			uriString = "spotify:album:" + id
		} else if strings.Contains(*recordUrl, "/playlist/") {
			uriString = "spotify:playlist:" + id
		} else {
			log.Fatalf("URL represents neither album nor playlist: %s", *recordUrl)
		}

		uri = spotify.URI(uriString)

		fmt.Println(uri)
	} else {
		if *deviceName == "" {
			log.Fatal("Device name is required.")
		}

		var spotifyUri spotify.URI

		if *playUri != "" && *readPath != "" {
			log.Fatal("Need to specify either a path or a URI, but not both.")
		} else if *playUri != "" {
			spotifyUri = spotify.URI(*playUri)
		} else {
			spotifyUri = getPlayURI(*readPath)
		}

		devices, err := client.PlayerDevices()
		handleError(err)

		var playerId *spotify.ID

		for _, device := range devices {
			if device.Name == *deviceName {
				playerId = &device.ID
				if !device.Active {
					err = client.TransferPlayback(*playerId, false)
					handleError(err)
				}
			}
		}

		if playerId == nil {
			log.Fatal("Player not found.")
		}

		playOptions := &spotify.PlayOptions{
			DeviceID:        playerId,
			PlaybackContext: &spotifyUri,
		}

		err = client.PlayOpt(playOptions)
		handleError(err)
	}
}

func handleError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func completeAuth(w http.ResponseWriter, r *http.Request) {
	tok, err := auth.Token(state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, state)
	}

	saveToken("tokenFile", tok)

	// use the token to get an authenticated client
	client := auth.NewClient(tok)
	fmt.Fprintf(w, "Login Completed!")
	ch <- &client
}

// Retrieves a token from a local file.
func tokenFromFile(filePath string) (*oauth2.Token, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	token := &oauth2.Token{}
	err = json.NewDecoder(file).Decode(token)
	return token, err
}

func getPlayURI(filePath string) spotify.URI {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var uri spotify.URI

	for scanner.Scan() {
		uri = spotify.URI(scanner.Text())
		break
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	} else if string(uri) == "" {
		log.Fatal("file empty or invalid")
	}

	return uri
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer file.Close()
	json.NewEncoder(file).Encode(token)
}
