package utils

import (
	"github.com/spf13/viper"
)

// Config stores all configuration of the application.
// The values are read by viper from a config file or environment variable.
type Config struct {
	DiscordToken    string `mapstructure:"DISCORD_TOKEN"`
	BaseURL         string `mapstructure:"BASE_URL"`
	ServerAddress   string `mapstructure:"SERVER_ADDRESS"`
	RethinkURL      string `mapstructure:"RETHINK_URL"`
	RethinkDatabase string `mapstructure:"RETHINK_DATABASE"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
