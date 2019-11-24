package diskplayer

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
	"net/url"
	"os"
	"time"
)

// NewAuthenticator returns a Spotify authenticator object configured with the required callback URL,
// client IT and client secret. An error is returned if one is encountered
func NewAuthenticator() (*spotify.Authenticator, error) {
	r := ConfigValue(SPOTIFY_CALLBACK_URL)
	u, err := url.Parse(r)
	if err != nil {
		return nil, err
	}

	id := ConfigValue(SPOTIFY_CLIENT_ID)
	s := ConfigValue(SPOTIFY_CLIENT_SECRET)

	auth := spotify.NewAuthenticator(u.String(), spotify.ScopeUserReadPrivate, spotify.ScopePlaylistReadPrivate,
		spotify.ScopeUserModifyPlaybackState, spotify.ScopeUserReadPlaybackState)

	auth.SetAuthInfo(id, s)

	return &auth, nil
}

// NewToken will create a new OAuth2 token request.
// The user will be prompted to visit a URL, and after access is granted a new OAuth2 token is returned.
// An error is returned if encountered.
func NewToken(ds DiskplayerServer) (*oauth2.Token, error) {
	s, err := ds.RunCallbackServer()
	if err != nil {
		return nil, err
	}

	u := ds.Authenticator().AuthURL(STATE_IDENTIFIER)
	fmt.Println("Please log in to Spotify by visiting the following page in your browser:", u)

	t := <-ds.TokenChannel()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = s.Shutdown(ctx)
	if err != nil {
		return nil, err
	}

	return t, nil
}

// ReadToken will attempt to deserialize a token whose path is defined in the diskplayer.yaml
// configuration file under the token.file_path field.
// Returns a pointer to an oauth2 token object or any error encountered.
func ReadToken() (*oauth2.Token, error) {
	p := ConfigValue(TOKEN_PATH)

	f, err := os.Open(p)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	t := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(t)
	return t, err
}

// SaveToken will serialize the provided token and save it to the file whose path is defined in the diskplayer. yaml
// configuration file under the token.file_path field.
// Returns an error if one is encountered.
func SaveToken(token *oauth2.Token) error {
	p := ConfigValue(TOKEN_PATH)

	f, err := os.OpenFile(p, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(token)
}
