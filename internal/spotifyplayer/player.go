package spotifyplayer

import (
	"context"
	"fmt"
	"github.com/dinofizz/diskplayer/internal/email"
	"github.com/zmb3/spotify"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	auth = spotify.NewAuthenticator(redirectURI, spotify.ScopeUserReadPrivate, spotify.ScopePlaylistReadPrivate,
		spotify.ScopeUserModifyPlaybackState, spotify.ScopeUserReadPlaybackState)
	ch          = make(chan *spotify.Client, 1)
	state       = "abc123"
	redirectURI = os.Getenv(SPOTIFY_CALLBACK_URL)
)

const (
	SPOTIFY_CALLBACK_URL = "SPOTIFY_CALLBACK_URL"
	SPOTIFY_DEVICE_NAME  = "SPOTIFY_DEVICE_NAME"
)

func Play(uri string) {
	deviceName := os.Getenv(SPOTIFY_DEVICE_NAME)
	if deviceName == "" {
		log.Fatalf("Environment variable %s is empty.", SPOTIFY_DEVICE_NAME)
	}

	if uri == "" {
		log.Fatal("Spotify URI is required.")
	}

	spotifyUri := spotify.URI(uri)

	client := client()

	playerId := getPlayerId(client, deviceName)

	playOptions := &spotify.PlayOptions{
		DeviceID:        playerId,
		PlaybackContext: &spotifyUri,
	}

	err := client.PlayOpt(playOptions)
	handleError(err)
}

func Pause() {
	err := client().Pause()
	handleError(err)
}

func client() *spotify.Client {
	http.HandleFunc("/callback", completeAuth)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got request for:", r.URL.String())
	})

	var server *http.Server

	token, err := tokenFromFile("tokenFile")
	if err != nil {
		if redirectURI == "" {
			log.Fatalf("Environment variable %s is empty.", SPOTIFY_CALLBACK_URL)
		}
		server = &http.Server{Addr: ":8080", Handler: nil}
		go server.ListenAndServe()
		url := auth.AuthURL(state)
		fmt.Println("Please log in to Spotify by visiting the following page in your browser:", url)
		_, err := email.SendAuthenticationUrlEmail(url)
		handleError(err)
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

func getPlayerId(client *spotify.Client, deviceName string) *spotify.ID {
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

	return playerId
}

func handleError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}
