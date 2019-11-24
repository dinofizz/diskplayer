package diskplayer

import (
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

// Returns an authenticated Spotify client object, or an error if encountered.
func NewClient(a *spotify.Authenticator, t *oauth2.Token) *SpotifyClient {
	c := a.NewClient(t)
	return &SpotifyClient{client: &c}
}

type Client interface {
	PlayerDevices() ([]spotify.PlayerDevice, error)
	Pause() error
	TransferPlayback(deviceID spotify.ID, play bool) error
	PlayOpt(opt *spotify.PlayOptions) error
}

type SpotifyClient struct {
	client *spotify.Client
}

func (sc *SpotifyClient) PlayerDevices() ([]spotify.PlayerDevice, error) {
	return sc.client.PlayerDevices()
}

func (sc *SpotifyClient) Pause() error {
	return sc.client.Pause()
}

func (sc *SpotifyClient) TransferPlayback(deviceID spotify.ID, play bool) error {
	return sc.client.TransferPlayback(deviceID, play)
}

func (sc *SpotifyClient) PlayOpt(opt *spotify.PlayOptions) error {
	return sc.client.PlayOpt(opt)
}
