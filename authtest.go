// This example demonstrates how to authenticate with Spotify using the authorization code flow.
// In order to run this example yourself, you'll need to:
//
//  1. Register an application at: https://developer.spotify.com/my-applications/
//       - Use "http://localhost:8080/callback" as the redirect URI
//  2. Set the SPOTIFY_ID environment variable to the client ID you got in step 1.
//  3. Set the SPOTIFY_SECRET environment variable to the client secret from step 1.
package main

import (
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/zmb3/spotify"
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
	// first start an HTTP server
	http.HandleFunc("/callback", completeAuth)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got request for:", r.URL.String())
	})
	go http.ListenAndServe(":8080", nil)

	token, err := tokenFromFile("tokenFile")
	if err != nil {
		url := auth.AuthURL(state)
		fmt.Println("Please log in to Spotify by visiting the following page in your browser:", url)
	} else {
		fmt.Println("Found an existing token file!")
		client := auth.NewClient(token)
		fmt.Println("Sending client to channel")
		ch <- &client
	}

	fmt.Println("Waiting for auth token...")
	// wait for auth to complete
	client := <-ch

	// use the client to make calls that require authorization
	user, err := client.CurrentUser()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("You are logged in as:", user.ID)

	page, err := client.GetPlaylistsForUser(user.ID)
	if err != nil {
		log.Fatal(err)
	}

	var starredID spotify.ID

	for _, playlist := range page.Playlists {
		fmt.Println(playlist.Name)
		if playlist.Name == "Starred" {
			starredID = playlist.ID
		}

	}

	devices, err := client.PlayerDevices()
	if err != nil {
		log.Fatal(err)
	}

	var activeId *spotify.ID

	for _, device := range devices {
		fmt.Println(device.Name, device.ID, device.Active)
		if device.Name == "Spotifyd@dell7490" {
			activeId = &device.ID
		}
	}

	if activeId == nil {
		log.Fatal("No active device")
	}

	err = client.TransferPlayback(*activeId, true)
	if err != nil {
		log.Fatal(err)
	}

	fullPlaylist, err := client.GetPlaylist(starredID)
	if err != nil {
		log.Fatal(err)
	}

	//totalTracks := fullPlaylist.Tracks.Total
	//fmt.Println(totalTracks)
	//trackNum := rand.Intn(int(totalTracks))
	rand.Seed(time.Now().UTC().UnixNano())
	trackNum := rand.Intn(100)
	fmt.Println(trackNum)

	uri := fullPlaylist.Tracks.Tracks[trackNum].Track.URI
	track := fullPlaylist.Tracks.Tracks[trackNum].Track

	uriList := []spotify.URI{uri}

	playOptions := &spotify.PlayOptions{
		DeviceID: activeId,
		URIs:     uriList,
	}

	err = client.PlayOpt(playOptions)
	if err != nil {
		log.Fatal(err)
	}

	fullTrack, err := client.GetTrack(track.ID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(fullTrack.Name)
	fmt.Println(fullTrack.Artists[0].Name)
	fmt.Println(fullTrack)
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
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}
