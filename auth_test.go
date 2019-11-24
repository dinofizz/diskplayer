package diskplayer

import (
	"errors"
	"github.com/dinofizz/diskplayer/mocks"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"testing"
	"time"
)

func TestNewAuthenticator(t *testing.T) {
	viper.Set("spotify.callback_url", "http://localhost:8732/callback")
	viper.Set("spotify.client_id", "client_id")
	viper.Set("spotify.client_secret", "client_secret")
	_, err := NewAuthenticator()
	assert.NoError(t, err)
}

func TestNewAuthenticatorParseURLError(t *testing.T) {
	viper.Set("spotify.callback_url", "//x:!*)florble")
	_, err := NewAuthenticator()
	assert.EqualError(t, err, "parse //x:!*)florble: invalid port \":!*)florble\" after host")
}

func TestNewToken(t *testing.T) {
	viper.Set("spotify.callback_url", "http://localhost:8732/callback")
	viper.Set("spotify.client_id", "client_id")
	viper.Set("spotify.client_secret", "client_secret")

	ms := new(mocks.DiskplayerServer)

	s := &http.Server{}

	a, err := NewAuthenticator()
	assert.NoError(t, err)

	ch := make(chan *oauth2.Token, 1)
	ms.On("TokenChannel").Return(ch)
	ms.On("Authenticator").Return(a, nil)
	ms.On("RunCallbackServer").Return(s, nil)

	var wg sync.WaitGroup
	wg.Add(1)

	tokExpected := oauth2.Token{
		AccessToken:  "temp_access_token",
		TokenType:    "Bearer",
		RefreshToken: "temp_refresh_token",
		Expiry:       time.Time{},
	}

	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		tokActual, e := NewToken(ms)
		assert.NoError(t, e)
		assert.Equal(t, tokActual.TokenType, tokExpected.TokenType)
		assert.Equal(t, tokActual.AccessToken, tokExpected.AccessToken)
		assert.Equal(t, tokActual.RefreshToken, tokExpected.RefreshToken)
		assert.Equal(t, tokActual.Expiry, tokExpected.Expiry)
	}(&wg)

	ch <- &tokExpected
	wg.Wait()
}

func TestNewTokenCallbackServerError(t *testing.T) {
	viper.Set("spotify.callback_url", "http://localhost:8732/callback")
	viper.Set("spotify.client_id", "client_id")
	viper.Set("spotify.client_secret", "client_secret")

	ms := new(mocks.DiskplayerServer)

	a, err := NewAuthenticator()
	assert.NoError(t, err)

	ms.On("Authenticator").Return(a, nil)
	ms.On("RunCallbackServer").Return(nil, errors.New("RunCallbackServer error"))
	tok, err := NewToken(ms)
	assert.Nil(t, tok)
	assert.EqualError(t, err, "RunCallbackServer error")
}

func TestReadToken(t *testing.T) {
	viper.Set("token.path", "./test-fixtures/test_token.json")
	tok, err := ReadToken()
	assert.NoError(t, err)
	assert.NotNil(t, tok)
	assert.Equal(t, "test_access_token", tok.AccessToken)
	assert.Equal(t, "test_refresh_token", tok.RefreshToken)
	assert.Equal(t, "2019-11-09 14:01:07.170612451 +0000 UTC", tok.Expiry.String())
	assert.Equal(t, "Bearer", tok.TokenType)
}

func TestReadTokenFileError(t *testing.T) {
	viper.Set("token.path", "./test-fixtures/not_a_real_path.json")
	tok, err := ReadToken()
	assert.Nil(t, tok)
	assert.Error(t, err)
	_, ok := err.(*os.PathError)
	assert.True(t, ok)
}

func TestSaveToken(t *testing.T) {
	const p = "./test-fixtures/temp_test_token.json"
	viper.Set("token.path", p)
	if _, err := os.Stat(p); os.IsExist(err) {
		err := os.Remove(p)
		if err != nil {
			t.Fatalf("Remove %s: %v", p, err)
		}
	}

	tok := oauth2.Token{
		AccessToken:  "temp_access_token",
		TokenType:    "Bearer",
		RefreshToken: "temp_refresh_token",
		Expiry:       time.Time{},
	}

	err := SaveToken(&tok)
	assert.NoError(t, err)

	err = os.Remove(p)
	assert.NoErrorf(t, err, "Failed to remove temporary test file: %s", p)
}

func TestSaveTokenWriteError(t *testing.T) {
	const p = "./test-fixtures/temp_test_token.json"
	viper.Set("token.path", p)
	f, err := os.Stat(p)
	if f != nil || os.IsExist(err) {
		err := os.Remove(p)
		if err != nil {
			t.Fatalf("Remove %s: %v", p, err)
		}
	}

	shmorp := []byte("shmorp")
	err = ioutil.WriteFile(p, shmorp, 0444)
	if err != nil {
		t.Fatalf("WriteFile %s: %v", p, err)
	}
	defer func() {
		err := os.Remove(p)
		assert.NoErrorf(t, err, "Failed to remove temporary test file: %s", p)
	}()

	tok := oauth2.Token{
		AccessToken:  "temp_access_token",
		TokenType:    "Bearer",
		RefreshToken: "temp_refresh_token",
		Expiry:       time.Time{},
	}

	err = SaveToken(&tok)
	assert.Error(t, err)
	assert.Equal(t, "open ./test-fixtures/temp_test_token.json: permission denied", err.Error())
}
