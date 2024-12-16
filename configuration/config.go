package configuration

import (
	"sync"

	"github.com/spf13/viper"
)

type Config struct {
	HTTPServerAddress string            `mapstructure:"HTTP_SERVER_ADDRESS"`
	DSN               string            `mapstructure:"DSN"`
	TokenSymmetricKey string            `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	RedisAddress      string            `mapstructure:"REDIS_ADDRESS"`
	RedisPassword     string            `mapstructure:"REDIS_PASSWORD"`
	RedisDB           int               `mapstructure:"REDIS_DB"`
	GoogleOAuthConfig GoogleOAuthConfig `mapstructure:",squash"`
}

type GoogleOAuthConfig struct {
	RedirectURL  string `mapstructure:"GOOGLE_OAUTH_CALL_BACK_URL"`
	ClientID     string `mapstructure:"GOOGLE_OAUTH_CLIENT_ID"`
	ClientSecret string `mapstructure:"GOOGLE_OAUTH_CLIENT_SECRET"`
}

var (
	once   sync.Once
	config *Config
)

func LoadConfig(path string) (config Config, err error) {
	viper.SetConfigName("app")
	viper.AddConfigPath(path)
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()

	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
