package diskplayer

import (
	"bufio"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

var recordTests = []struct {
	in  string
	out string
	e   string
}{
	{"https://open.spotify.com/album/1S7mumn7D4riEX2gVWYgPO", "spotify:album:1S7mumn7D4riEX2gVWYgPO", ""},
	{"https://open.spotify.com/playlist/5XsXwH5uWdhpAWsigjWMTA", "spotify:playlist:5XsXwH5uWdhpAWsigjWMTA", ""},
	{"florble", "",
		"URL represents neither album nor playlist: florble"},
}

func TestRecord(t *testing.T) {
	const p = "./test_recorder_path.contents"
	viper.Set("recorder.file_path", p)

	for _, tt := range recordTests {
		t.Run(tt.in, func(t *testing.T) {

			err := Record(tt.in, p)
			if tt.e == "" {
				assert.NoErrorf(t, err, "Record encountered an unexpected error.")
			} else if err != nil && tt.e != "" {
				assert.Equal(t, tt.e, err.Error(), "")
				return
			} else {
				t.Error("Expected an error but did not receive one.")
			}

			defer func() {
				err := os.Remove(p)
				assert.NoErrorf(t, err, "Failed to remove temporary test file: %s", p)
			}()
			f, err := os.Open(p)
			assert.NoErrorf(t, err, "Failed to open test path: %s", p)

			defer f.Close()

			s := bufio.NewScanner(f)
			var l string
			var n int
			for s.Scan() {
				l = s.Text()
				n++
			}

			assert.Equal(t, 1, n, "Expected to only read one line from recorder path")
			assert.Equal(t, tt.out, l, "Recorded Spotify URI was not as expected.")
		})
	}

}

func TestRecordWriteError(t *testing.T) {
	const p = "./test_recorder_path.contents"
	viper.Set("recorder.file_path", p)
	if _, err := os.Stat(p); os.IsExist(err) {
		err := os.Remove(p)
		if err != nil {
			t.Fatalf("Remove %s: %v", p, err)
		}
	}

	shmorp := []byte("shmorp")
	err := ioutil.WriteFile(p, shmorp, 0444)
	if err != nil {
		t.Fatalf("WriteFile %s: %v", p, err)
	}
	defer func() {
		err := os.Remove(p)
		assert.NoErrorf(t, err, "Failed to remove temporary test file: %s", p)
	}()
	err = Record(recordTests[0].in, p)
	assert.Error(t, err)
	assert.Equal(t, "open ./test_recorder_path.contents: permission denied", err.Error())
}
