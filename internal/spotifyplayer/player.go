package spotifyplayer

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dinofizz/diskplayer/internal/config"
	"github.com/dinofizz/diskplayer/internal/email"
	"github.com/dinofizz/diskplayer/internal/errorhandler"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	ch    = make(chan *spotify.Client, 1)
	state = "abc123"
	auth  spotify.Authenticator
)

const (
	SPOTIFY_CALLBACK_URL   = "spotify.callback_url"
	SPOTIFY_DEVICE_NAME    = "spotify.device_name"
	SPOTIFY_CLIENT_ID      = "spotify.client_id"
	SPOTIFY_CLIENT_SECRET  = "spotify.client_secret"
	SPOTIFY_ID_ENV_VAR     = "SPOTIFY_ID"
	SPOTIFY_SECRET_ENV_VAR = "SPOTIFY_SECRET"
	TOKEN_PATH             = "token.path"
)

func Play(uri string) {
	deviceName := config.GetConfigString(SPOTIFY_DEVICE_NAME)

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
	errorhandler.HandleError(err)
}

func Pause() {
	err := client().Pause()
	errorhandler.HandleError(err)
}

func client() *spotify.Client {

	var server *http.Server

	token, err := tokenFromFile()
	if err != nil {
		server = fetchNewToken()
	} else {
		newAuthenticator()
		client := auth.NewClient(token)
		ch <- &client
	}

	client := <-ch

	if server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := server.Shutdown(ctx)
		errorhandler.HandleError(err)
	}

	return client
}

func fetchNewToken() *http.Server {
	newAuthenticator()
	http.HandleFunc("/callback", completeAuth)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got request for:", r.URL.String())
	})
	server := &http.Server{Addr: ":8080", Handler: nil}
	go server.ListenAndServe()
	url := auth.AuthURL(state)
	fmt.Println("Please log in to Spotify by visiting the following page in your browser:", url)
	_, err := email.SendAuthenticationUrlEmail(url)
	errorhandler.HandleError(err)
	return server
}

func newAuthenticator() {
	redirectURI := config.GetConfigString(SPOTIFY_CALLBACK_URL)
	clientId := config.GetConfigString(SPOTIFY_CLIENT_ID)
	clientSecret := config.GetConfigString(SPOTIFY_CLIENT_SECRET)

	// Unset any existing environment variables
	err := os.Unsetenv(SPOTIFY_ID_ENV_VAR)
	errorhandler.HandleError(err)
	err = os.Unsetenv(SPOTIFY_SECRET_ENV_VAR)
	errorhandler.HandleError(err)

	// Set the environment variables required for Spotify auth
	err = os.Setenv(SPOTIFY_ID_ENV_VAR, clientId)
	errorhandler.HandleError(err)
	err = os.Setenv(SPOTIFY_SECRET_ENV_VAR, clientSecret)
	errorhandler.HandleError(err)

	auth = spotify.NewAuthenticator(redirectURI, spotify.ScopeUserReadPrivate, spotify.ScopePlaylistReadPrivate,
		spotify.ScopeUserModifyPlaybackState, spotify.ScopeUserReadPlaybackState)
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

	saveToken(tok)

	client := auth.NewClient(tok)
	fmt.Fprintf(w, "Login Completed!")
	ch <- &client
}

func getPlayerId(client *spotify.Client, deviceName string) *spotify.ID {
	devices, err := client.PlayerDevices()
	errorhandler.HandleError(err)

	var playerId *spotify.ID
	for _, device := range devices {
		if device.Name == deviceName {
			playerId = &device.ID
			if !device.Active {
				err = client.TransferPlayback(*playerId, false)
				errorhandler.HandleError(err)
			}
		}
	}

	if playerId == nil {
		log.Fatal("Player not found.")
	}

	return playerId
}

// Retrieves a token from a local file.
func tokenFromFile() (*oauth2.Token, error) {
	tokenpath := config.GetConfigString(TOKEN_PATH)

	file, err := os.Open(tokenpath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	token := &oauth2.Token{}
	err = json.NewDecoder(file).Decode(token)
	return token, err
}

// Saves a token to a file path.
func saveToken(token *oauth2.Token) {
	tokenpath := config.GetConfigString(TOKEN_PATH)

	file, err := os.OpenFile(tokenpath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer file.Close()
	json.NewEncoder(file).Encode(token)
}
