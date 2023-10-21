package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	SaltLength int
	PathUpload string
}

func ProvideConfig() Config {
	return Config{
		SaltLength: viper.GetInt("SALT_LENGTH"),
		PathUpload: viper.GetString("PATH_UPLOAD"),
	}
}
