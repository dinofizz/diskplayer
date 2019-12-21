package diskplayer

import (
	"errors"
	"fmt"
	"github.com/dinofizz/diskplayer/mocks"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/zmb3/spotify"
	"os"
	"testing"
)

func TestPlayUriSuccess(t *testing.T) {
	viper.Set("spotify.device_name", "test_device_name")

	m := new(mocks.Client)

	d := spotify.PlayerDevice{
		ID:         "TEST_ID",
		Active:     false,
		Restricted: false,
		Name:       "test_device_name",
		Type:       "",
		Volume:     0,
	}

	ds := []spotify.PlayerDevice{d}

	m.On("PlayerDevices").Return(ds, nil)
	m.On("PlayOpt", mock.AnythingOfType("*spotify.PlayOptions")).Return(nil)

	err := PlayUri(m, "foobar")
	assert.NoError(t, err)
}

func TestPlayPathSuccess(t *testing.T) {
	viper.Set("spotify.device_name", "test_device_name")

	m := new(mocks.Client)

	d := spotify.PlayerDevice{
		ID:         "TEST_ID",
		Active:     false,
		Restricted: false,
		Name:       "test_device_name",
		Type:       "",
		Volume:     0,
	}

	ds := []spotify.PlayerDevice{d}

	m.On("PlayerDevices").Return(ds, nil)
	m.On("PlayOpt", mock.AnythingOfType("*spotify.PlayOptions")).Return(nil)

	err := PlayPath(m, "./test-fixtures/diskplayer.contents")
	assert.NoError(t, err)
}

func TestPlayPathInvalidPath(t *testing.T) {
	m := new(mocks.Client)
	err := PlayPath(m, "./test-fixtures/not_a_real_path")
	if _, ok := err.(*os.PathError); !ok {
		t.Errorf("PlayPath was expected to fail with os.PathError.")
	}
}

func TestPlayPathEmptyFileError(t *testing.T) {
	m := new(mocks.Client)
	err := PlayPath(m, "./test-fixtures/empty.contents")
	assert.EqualError(t, err, "unable to read line from path: ./test-fixtures/empty.contents")
}

func TestPlayUriNoUriError(t *testing.T) {
	m := new(mocks.Client)
	err := PlayUri(m, "")
	assert.EqualError(t, err, "spotify URI is required")
}

func TestPlayUrPlayerDevicesError(t *testing.T) {
	viper.Set("spotify.device_name", "test_device_name")

	m := new(mocks.Client)
	var ds []spotify.PlayerDevice

	const e string = "PlayerDevices error"

	m.On("PlayerDevices").Return(ds, errors.New(e))
	err := PlayUri(m, "dummy_uri")
	assert.EqualError(t, err, e)
}

func TestPlayUriDeviceNotFoundError(t *testing.T) {
	const n = "test_device_name"
	viper.Set("spotify.device_name", n)

	m := new(mocks.Client)

	d := spotify.PlayerDevice{
		ID:         "TEST_ID",
		Active:     false,
		Restricted: false,
		Name:       "another_device_name",
		Type:       "",
		Volume:     0,
	}

	ds := []spotify.PlayerDevice{d}

	m.On("PlayerDevices").Return(ds, nil)

	err := PlayUri(m, "foobar")
	assert.EqualError(t, err, fmt.Sprintf("client identified by %s not found", n))
}

func TestPlayUriTransferPlaybackSuccess(t *testing.T) {
	const n = "test_device_name"
	viper.Set("spotify.device_name", n)

	m := new(mocks.Client)

	ds := []spotify.PlayerDevice{
		{
			ID:         "ANOTHER_TEST_ID",
			Active:     true,
			Restricted: false,
			Name:       "another_device_name",
			Type:       "",
			Volume:     0,
		},
		{
			ID:         "TEST_ID",
			Active:     false,
			Restricted: false,
			Name:       n,
			Type:       "",
			Volume:     0,
		},
	}

	m.On("PlayerDevices").Return(ds, nil)
	m.On("Pause").Return(nil)
	m.On("TransferPlayback", mock.AnythingOfType("spotify.ID"), false).Return(nil)
	m.On("PlayOpt", mock.AnythingOfType("*spotify.PlayOptions")).Return(nil)

	err := PlayUri(m, "foobar")
	assert.NoError(t, err)
}

func TestPlayUriTransferPlaybackPauseError(t *testing.T) {
	const n = "test_device_name"
	viper.Set("spotify.device_name", n)

	m := new(mocks.Client)

	ds := []spotify.PlayerDevice{
		{
			ID:         "ANOTHER_TEST_ID",
			Active:     true,
			Restricted: false,
			Name:       "another_device_name",
			Type:       "",
			Volume:     0,
		},
		{
			ID:         "TEST_ID",
			Active:     false,
			Restricted: false,
			Name:       n,
			Type:       "",
			Volume:     0,
		},
	}

	m.On("PlayerDevices").Return(ds, nil)
	const e = "pause error"
	m.On("Pause").Return(errors.New(e))

	err := PlayUri(m, "foobar")
	assert.EqualError(t, err, e)
}

func TestPlayUriTransferPlaybackTransferError(t *testing.T) {
	const n = "test_device_name"
	viper.Set("spotify.device_name", n)

	m := new(mocks.Client)

	ds := []spotify.PlayerDevice{
		{
			ID:         "ANOTHER_TEST_ID",
			Active:     true,
			Restricted: false,
			Name:       "another_device_name",
			Type:       "",
			Volume:     0,
		},
		{
			ID:         "TEST_ID",
			Active:     false,
			Restricted: false,
			Name:       n,
			Type:       "",
			Volume:     0,
		},
	}

	m.On("PlayerDevices").Return(ds, nil)
	const e = "transfer error"
	m.On("Pause").Return(nil)
	m.On("TransferPlayback", mock.AnythingOfType("spotify.ID"), false).Return(errors.New(e))
	m.On("PlayOpt", mock.AnythingOfType("*spotify.PlayOptions")).Return(nil)

	err := PlayUri(m, "foobar")
	assert.EqualError(t, err, e)
}

func TestPauseSuccess(t *testing.T) {
	const n = "test_device_name"
	viper.Set("spotify.device_name", n)

	m := new(mocks.Client)

	ds := []spotify.PlayerDevice{
		{
			ID:         "ANOTHER_TEST_ID",
			Active:     true,
			Restricted: false,
			Name:       "another_device_name",
			Type:       "",
			Volume:     0,
		},
		{
			ID:         "TEST_ID",
			Active:     false,
			Restricted: false,
			Name:       n,
			Type:       "",
			Volume:     0,
		},
	}

	m.On("PlayerDevices").Return(ds, nil)
	m.On("Pause").Return(nil)

	err := Pause(m)
	assert.NoError(t, err)
}

func TestPauseNoneActiveSuccess(t *testing.T) {
	const n = "test_device_name"
	viper.Set("spotify.device_name", n)

	m := new(mocks.Client)

	ds := []spotify.PlayerDevice{
		{
			ID:         "ANOTHER_TEST_ID",
			Active:     false,
			Restricted: false,
			Name:       "another_device_name",
			Type:       "",
			Volume:     0,
		},
		{
			ID:         "TEST_ID",
			Active:     false,
			Restricted: false,
			Name:       n,
			Type:       "",
			Volume:     0,
		},
	}

	m.On("PlayerDevices").Return(ds, nil)
	m.On("Pause").Return(nil)

	err := Pause(m)
	assert.NoError(t, err)
}

func TestPauseDeviceNotFound(t *testing.T) {
	const n = "test_device_name"
	viper.Set("spotify.device_name", n)

	m := new(mocks.Client)

	ds := []spotify.PlayerDevice{
		{
			ID:         "ANOTHER_TEST_ID",
			Active:     true,
			Restricted: false,
			Name:       "another_device_name",
			Type:       "",
			Volume:     0,
		},
	}

	m.On("PlayerDevices").Return(ds, nil)

	err := Pause(m)
	assert.EqualError(t, err, fmt.Sprintf("client identified by %s not found", n))
}

func TestPausePlayerDevicesError(t *testing.T) {
	viper.Set("spotify.device_name", "test_device_name")

	m := new(mocks.Client)
	var ds []spotify.PlayerDevice

	const e string = "PlayerDevices error"

	m.On("PlayerDevices").Return(ds, errors.New(e))
	err := Pause(m)
	assert.EqualError(t, err, e)
}

func TestPauseError(t *testing.T) {
	const n = "test_device_name"
	viper.Set("spotify.device_name", n)

	m := new(mocks.Client)

	ds := []spotify.PlayerDevice{
		{
			ID:         "TEST_ID",
			Active:     true,
			Restricted: false,
			Name:       n,
			Type:       "",
			Volume:     0,
		},
	}

	m.On("PlayerDevices").Return(ds, nil)
	const e = "pause error"
	m.On("Pause").Return(errors.New(e))

	err := Pause(m)
	assert.EqualError(t, err, e)
}
