package config

import (
	"github.com/ory/viper"
)

type Config struct {
	SaltLength int
}

func ProvideConfig() Config {
	return Config{
		SaltLength: viper.GetInt("SALT_LENGTH"),
	}
}
