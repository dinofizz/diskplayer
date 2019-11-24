package diskplayer

import (
	"context"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
	"testing"
	"time"
)

func TestNewDiskplayerServer(t *testing.T) {
	a := &spotify.Authenticator{}
	ch := make(chan *oauth2.Token, 1)
	s := NewDiskplayerServer(a, ch)
	assert.Equal(t, a, s.Authenticator())
	assert.Equal(t, ch, s.TokenChannel())
}

func TestRealDiskplayerServer_RunCallbackServer(t *testing.T) {
	viper.Set("spotify.callback_url", "http://localhost:8732/callback")
	a := &spotify.Authenticator{}
	ch := make(chan *oauth2.Token, 1)
	ds := NewDiskplayerServer(a, ch)
	viper.Set("recorder.server_port", 4389)
	s, err := ds.RunCallbackServer()
	assert.NoError(t,err)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = s.Shutdown(ctx)
	assert.NoError(t,err)
}
