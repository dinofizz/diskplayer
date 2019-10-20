package spotifyplayer

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	auth = spotify.NewAuthenticator(redirectURI, spotify.ScopeUserReadPrivate, spotify.ScopePlaylistReadPrivate,
		spotify.ScopeUserModifyPlaybackState, spotify.ScopeUserReadPlaybackState)
	ch    = make(chan *spotify.Client, 1)
	state = "abc123"
)

const redirectURI = "http://localhost:8080/callback"

func Play(deviceName string, uriString string) {
	if deviceName == "" {
		log.Fatal("Device name is required.")
	}

	if uriString == "" {
		log.Fatal("Spotify URI is required.")
	}

	spotifyUri := spotify.URI(uriString)

	client := getClient()

	devices, err := client.PlayerDevices()
	handleError(err)

	var playerId *spotify.ID

	for _, device := range devices {
		if device.Name == deviceName {
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

func Pause() {
	err := getClient().Pause()
	handleError(err)
	os.Exit(0)
}

func getClient() *spotify.Client {

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

	return client
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
