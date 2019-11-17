package diskplayer

import (
	"github.com/spf13/viper"
	"testing"
)

func TestReadConfigError(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	ReadConfig("unknown_path")
}

func TestConfigValue(t *testing.T) {
	viper.AddConfigPath("./test-fixtures")
	ReadConfig("test_config")

	if ConfigValue("spotify.device_name") != "my_device_name" {
		t.Error("Failed to retrieve value from test_config.yaml")
	}

}
