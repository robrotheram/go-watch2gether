package utils

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var Configuration = Config{}

// Config stores all configuration of the application.
// The values are read by viper from a config file or environment variable.
type Config struct {
	DiscordToken        string `mapstructure:"DISCORD_TOKEN"`
	DiscordClientID     string `mapstructure:"DISCORD_CLIENT_ID"`
	DiscordClientSecret string `mapstructure:"DISCORD_CLIENT_SECRET"`
	DiscordNotify       bool   `mapstructure:"DISCORD_ENABLE_NOTIFICATIONS"`
	SessionSecret       string `mapstructure:"SESSION_SECRET"`
	BaseURL             string `mapstructure:"BASE_URL"`
	ServerAddress       string `mapstructure:"SERVER_ADDRESS"`
	RethinkURL          string `mapstructure:"RETHINK_URL"`
	RethinkDatabase     string `mapstructure:"RETHINK_DATABASE"`
	Dev                 bool   `mapstructure:"DEVELOPMENT"`
	loglevel            string `mapstructure:"LOG_LEVEL"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&Configuration)
	return
}

func GetLoglevel() log.Level {
	switch Configuration.loglevel {
	case "fatal":
		return log.FatalLevel
	case "erro":
		return log.ErrorLevel
	case "warn":
		return log.WarnLevel
	case "info":
		return log.InfoLevel
	case "debug":
		return log.DebugLevel
	default:
		return log.InfoLevel
	}
}
