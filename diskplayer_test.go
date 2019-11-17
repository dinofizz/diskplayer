package diskplayer

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/mock"
	"github.com/zmb3/spotify"
	"os"
	"testing"
)

type MockSpotifyClient struct {
	mock.Mock
}

func (m *MockSpotifyClient) PlayerDevices() ([]spotify.PlayerDevice, error) {
	args := m.Called()
	return args.Get(0).([]spotify.PlayerDevice), args.Error(1)
}

func (m *MockSpotifyClient) Pause() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockSpotifyClient) TransferPlayback(deviceID spotify.ID, play bool) error {
	args := m.Called(deviceID, play)
	return args.Error(0)
}

func (m *MockSpotifyClient) PlayOpt(opt *spotify.PlayOptions) error {
	args := m.Called(opt)
	return args.Error(0)
}

func TestPlaySuccess(t *testing.T) {
	viper.Set("recorder.file_path", "./test-fixtures/diskplayer.contents")
	viper.Set("spotify.device_name", "test_device_name")

	m := new(MockSpotifyClient)

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

	err := Play(m)
	if err != nil {
		t.Errorf("Play failed with error %s", err)
	}
}
func TestPlayUriSuccess(t *testing.T) {
	viper.Set("spotify.device_name", "test_device_name")

	m := new(MockSpotifyClient)

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
	if err != nil {
		t.Errorf("PlayUri failed with error %s", err)
	}
}

func TestPlayPathSuccess(t *testing.T) {
	viper.Set("spotify.device_name", "test_device_name")

	m := new(MockSpotifyClient)

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
	if err != nil {
		t.Errorf("PlayPath failed with error %s", err)
	}
}

func TestPlayPathInvalidPath(t *testing.T) {
	m := new(MockSpotifyClient)
	err := PlayPath(m, "./test-fixtures/not_a_real_path")
	if _, ok := err.(*os.PathError); !ok {
		t.Errorf("PlayPath was expected to fail with os.PathError.")
	}
}

func TestPlayPathEmptyFileError(t *testing.T) {
	m := new(MockSpotifyClient)
	err := PlayPath(m, "./test-fixtures/empty.contents")
	if err == nil || err.Error() != "unable to read line from path: ./test-fixtures/empty.contents" {
		t.Error("PlayPath was expected to fail when reading an empty file.")
	}
}

func TestPlayUriNoUriError(t *testing.T) {
	m := new(MockSpotifyClient)
	err := PlayUri(m, "")
	if err == nil || err.Error() != "spotify URI is required" {
		t.Error("PlayUri was expected to fail when provided an empty URI.")
	}
}

func TestPlayUrPlayerDevicesError(t *testing.T) {
	viper.Set("spotify.device_name", "test_device_name")

	m := new(MockSpotifyClient)
	ds := []spotify.PlayerDevice{}

	const e string = "PlayerDevices error"

	m.On("PlayerDevices").Return(ds, errors.New(e))
	err := PlayUri(m, "dummy_uri")
	if err == nil || err.Error() != e {
		t.Error("PlayUri was expected to return error fetching PlayerDevices")
	}
}

func TestPlayUriDeviceNotFoundError(t *testing.T) {
	const n = "test_device_name"
	viper.Set("spotify.device_name", n)

	m := new(MockSpotifyClient)

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
	if err == nil || err.Error() != fmt.Sprintf("client identified by %s not found", n) {
		t.Error("PlayUri was expected to return error when the device name was not found.")
	}
}

func TestPlayUriTransferPlaybackSuccess(t *testing.T) {
	const n = "test_device_name"
	viper.Set("spotify.device_name", n)

	m := new(MockSpotifyClient)

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
	if err != nil {
		t.Errorf("PlayUri failed with error %s", err)
	}
}

func TestPlayUriTransferPlaybackPauseError(t *testing.T) {
	const n = "test_device_name"
	viper.Set("spotify.device_name", n)

	m := new(MockSpotifyClient)

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
	if err == nil || err.Error() != e {
		t.Error("PlayUri was expected to return error when pausing playback during playback transfer sequence.")
	}
}

func TestPlayUriTransferPlaybackTransferError(t *testing.T) {
	const n = "test_device_name"
	viper.Set("spotify.device_name", n)

	m := new(MockSpotifyClient)

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
	if err == nil || err.Error() != e {
		t.Error("PlayUri was expected to return error when transferring playback during playback transfer sequence.")
	}
}

func TestPauseSuccess(t *testing.T) {
	const n = "test_device_name"
	viper.Set("spotify.device_name", n)

	m := new(MockSpotifyClient)

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
	if err != nil {
		t.Errorf("Pause failed with error %s", err)
	}
}

func TestPauseNoneActiveSuccess(t *testing.T) {
	const n = "test_device_name"
	viper.Set("spotify.device_name", n)

	m := new(MockSpotifyClient)

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
	if err != nil {
		t.Errorf("Pause failed with error %s", err)
	}
}

func TestPauseDeviceNotFound(t *testing.T) {
	const n = "test_device_name"
	viper.Set("spotify.device_name", n)

	m := new(MockSpotifyClient)

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
	if err == nil || err.Error() != fmt.Sprintf("client identified by %s not found", n) {
		t.Error("Pause was expected to return error when the device name was not found.")
	}
}

func TestPausePlayerDevicesError(t *testing.T) {
	viper.Set("spotify.device_name", "test_device_name")

	m := new(MockSpotifyClient)
	ds := []spotify.PlayerDevice{}

	const e string = "PlayerDevices error"

	m.On("PlayerDevices").Return(ds, errors.New(e))
	err := Pause(m)
	if err == nil || err.Error() != e {
		t.Error("Pause was expected to return error fetching PlayerDevices")
	}
}
func TestPauseError(t *testing.T) {
	const n = "test_device_name"
	viper.Set("spotify.device_name", n)

	m := new(MockSpotifyClient)

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
	if err == nil || err.Error() != e {
		t.Error("Pause was expected to return error.")
	}
}
