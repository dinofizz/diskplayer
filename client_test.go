package diskplayer

import (
	"github.com/stretchr/testify/assert"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	tok := oauth2.Token{
		AccessToken:  "temp_access_token",
		TokenType:    "Bearer",
		RefreshToken: "temp_refresh_token",
		Expiry:       time.Time{},
	}

	a := spotify.Authenticator{}
	c := NewClient(&a, &tok)
	assert.NotNil(t, c)
}
