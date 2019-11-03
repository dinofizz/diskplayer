package spotifyplayer

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dinofizz/diskplayer/internal/config"
	"github.com/dinofizz/diskplayer/internal/email"
	"github.com/dinofizz/diskplayer/internal/errorhandler"
	"github.com/dinofizz/diskplayer/internal/server"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
	"io/ioutil"
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
	RECORD_PATH            = "recorder.file_path"
)

func Play() {
	p := config.GetConfigString(RECORD_PATH)
	u, err := ioutil.ReadFile(p)
	errorhandler.HandleError(err)
	PlayUri(string(u))
}

func PlayUri(u string) {
	if u == "" {
		log.Fatal("Spotify URI is required.")
	}

	spotifyUri := spotify.URI(u)

	c := client()

	d := config.GetConfigString(SPOTIFY_DEVICE_NAME)
	id := getPlayerId(c, d)

	o := &spotify.PlayOptions{
		DeviceID:        id,
		PlaybackContext: &spotifyUri,
	}

	err := c.PlayOpt(o)
	errorhandler.HandleError(err)
}

func Pause() {
	err := client().Pause()
	errorhandler.HandleError(err)
}

func client() *spotify.Client {

	var s *http.Server

	t, err := tokenFromFile()
	if err != nil {
		s = fetchNewToken()
	} else {
		newAuthenticator()
		c := auth.NewClient(t)
		ch <- &c
	}

	c := <-ch

	if s != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := s.Shutdown(ctx)
		errorhandler.HandleError(err)
	}

	return c
}

func fetchNewToken() *http.Server {
	newAuthenticator()
	s := server.RunCallbackServer(completeAuth)
	u := auth.AuthURL(state)
	fmt.Println("Please log in to Spotify by visiting the following page in your browser:", u)
	_, err := email.SendAuthenticationUrlEmail(u)
	errorhandler.HandleError(err)
	return s
}

func newAuthenticator() {
	r := config.GetConfigString(SPOTIFY_CALLBACK_URL)
	id := config.GetConfigString(SPOTIFY_CLIENT_ID)
	s := config.GetConfigString(SPOTIFY_CLIENT_SECRET)

	// Unset any existing environment variables
	err := os.Unsetenv(SPOTIFY_ID_ENV_VAR)
	errorhandler.HandleError(err)
	err = os.Unsetenv(SPOTIFY_SECRET_ENV_VAR)
	errorhandler.HandleError(err)

	// Set the environment variables required for Spotify auth
	err = os.Setenv(SPOTIFY_ID_ENV_VAR, id)
	errorhandler.HandleError(err)
	err = os.Setenv(SPOTIFY_SECRET_ENV_VAR, s)
	errorhandler.HandleError(err)

	auth = spotify.NewAuthenticator(r, spotify.ScopeUserReadPrivate, spotify.ScopePlaylistReadPrivate,
		spotify.ScopeUserModifyPlaybackState, spotify.ScopeUserReadPlaybackState)
}

func completeAuth(w http.ResponseWriter, r *http.Request) {
	t, err := auth.Token(state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, state)
	}

	saveToken(t)

	c := auth.NewClient(t)
	fmt.Fprintf(w, "Login Completed!")
	ch <- &c
}

func getPlayerId(c *spotify.Client, n string) *spotify.ID {
	ds, err := c.PlayerDevices()
	errorhandler.HandleError(err)

	var id *spotify.ID
	for _, d := range ds {
		if d.Name == n {
			id = &d.ID
			if !d.Active {
				err = c.TransferPlayback(*id, false)
				errorhandler.HandleError(err)
			}
		}
	}

	if id == nil {
		log.Fatal("Player not found.")
	}

	return id
}

// Retrieves a token from a local file.
func tokenFromFile() (*oauth2.Token, error) {
	p := config.GetConfigString(TOKEN_PATH)

	f, err := os.Open(p)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	t := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(t)
	return t, err
}

// Saves a token to a file path.
func saveToken(token *oauth2.Token) {
	p := config.GetConfigString(TOKEN_PATH)

	f, err := os.OpenFile(p, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}
