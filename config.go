package diskplayer

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
)

// ReadConfig reads in the configuration values from the diskplayer.yaml configuration file.
func ReadConfig() {
	viper.SetConfigName("diskplayer")
	viper.AddConfigPath("/etc/diskplayer/")
	viper.AddConfigPath("$HOME/.config/diskplayer/")
	viper.AddConfigPath(".")
	viper.SetDefault("token.path", "token.json")
	viper.SetDefault("spotify.callback_url", "http://localhost:8080/callback")
	viper.SetDefault("recorder.server_port", "3000")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}

// ConfigValue returns the configuration value identified by the provided key.
// If none is found the application quits with an error message and exit code 1.
func ConfigValue(key string) string {
	value := viper.GetString(key)
	if value == "" {
		log.Fatalf("Configuration value \"%s\" is empty.", key)
	}

	return value
}
