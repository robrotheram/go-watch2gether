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
	SessionSecret       string `mapstructure:"SESSION_SECRET"`
	BaseURL             string `mapstructure:"BASE_URL"`
	DatabasePath        string `mapstructure:"DATABASE_PATH"`
	Dev                 bool   `mapstructure:"DEVELOPMENT"`
	Reset               bool   `mapstructure:"RESET"`
	Loglevel            string `mapstructure:"LOG_LEVEL"`
	ListenPort          string `mapstructure:"LISTEN_PORT"`
	BetterStackToken    string `mapstructure:"BETTER_STACK_TOKEN"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) error {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	err = viper.Unmarshal(&Configuration)
	SetDefaults()
	return err
}

func (c *Config) GetLoglevel() log.Level {
	switch c.Loglevel {
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

func SetDefaults() {
	if Configuration.ListenPort == "" {
		Configuration.ListenPort = "8080"
	}
}
