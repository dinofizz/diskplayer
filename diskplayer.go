package diskplayer

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

func Play() {
	p := GetConfigString(RECORD_PATH)
	u, err := ioutil.ReadFile(p)
	handleError(err)
	PlayUri(string(u))
}

func PlayUri(u string) {
	if u == "" {
		log.Fatal("Spotify URI is required.")
	}

	spotifyUri := spotify.URI(u)

	c := client()

	d := GetConfigString(SPOTIFY_DEVICE_NAME)
	id := getPlayerId(c, d)

	o := &spotify.PlayOptions{
		DeviceID:        id,
		PlaybackContext: &spotifyUri,
	}

	err := c.PlayOpt(o)
	handleError(err)
}

func Pause() {
	err := client().Pause()
	handleError(err)
}

func client() *spotify.Client {

	var s *http.Server
	ch := make(chan *spotify.Client, 1)

	t, err := tokenFromFile()
	if err != nil {
		s = fetchNewToken(ch)
	} else {
		auth := newAuthenticator()
		c := auth.NewClient(t)
		ch <- &c
	}

	c := <-ch

	if s != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := s.Shutdown(ctx)
		handleError(err)
	}

	return c
}

func fetchNewToken(ch chan *spotify.Client) *http.Server {
	auth := newAuthenticator()
	h := CallbackHandler{ch: ch, auth: &auth}
	s := RunCallbackServer(h)
	u := auth.AuthURL(STATE_IDENTIFIER)
	fmt.Println("Please log in to Spotify by visiting the following page in your browser:", u)
	return s
}

func newAuthenticator() spotify.Authenticator {
	r := GetConfigString(SPOTIFY_CALLBACK_URL)
	u, err := url.Parse(r)
	handleError(err)

	id := GetConfigString(SPOTIFY_CLIENT_ID)
	s := GetConfigString(SPOTIFY_CLIENT_SECRET)

	// Unset any existing environment variables
	err = os.Unsetenv(SPOTIFY_ID_ENV_VAR)
	handleError(err)
	err = os.Unsetenv(SPOTIFY_SECRET_ENV_VAR)
	handleError(err)

	// Set the environment variables required for Spotify auth
	err = os.Setenv(SPOTIFY_ID_ENV_VAR, id)
	handleError(err)
	err = os.Setenv(SPOTIFY_SECRET_ENV_VAR, s)
	handleError(err)

	return spotify.NewAuthenticator(u.String(), spotify.ScopeUserReadPrivate, spotify.ScopePlaylistReadPrivate,
		spotify.ScopeUserModifyPlaybackState, spotify.ScopeUserReadPlaybackState)
}

type CallbackHandler struct {
	ch   chan *spotify.Client
	auth *spotify.Authenticator
}

func (h CallbackHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t, err := h.auth.Token(STATE_IDENTIFIER, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}
	if st := r.FormValue("state"); st != STATE_IDENTIFIER {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, STATE_IDENTIFIER)
	}

	saveToken(t)

	c := h.auth.NewClient(t)
	fmt.Fprintf(w, "Login Completed!")
	h.ch <- &c
}

func getPlayerId(c *spotify.Client, n string) *spotify.ID {
	ds, err := c.PlayerDevices()
	handleError(err)

	var id *spotify.ID
	for _, d := range ds {
		if d.Name == n {
			id = &d.ID
			if !d.Active {
				err := c.Pause()
				handleError(err)
				err = c.TransferPlayback(*id, false)
				handleError(err)
			}
			break
		}
	}

	if id == nil {
		log.Fatal("Player not found.")
	}

	return id
}

// Retrieves a token from a local file.
func tokenFromFile() (*oauth2.Token, error) {
	p := GetConfigString(TOKEN_PATH)

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
	p := GetConfigString(TOKEN_PATH)

	f, err := os.OpenFile(p, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func handleError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}
