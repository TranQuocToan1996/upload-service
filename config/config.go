package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	SaltLength       int
	PathUpload       string
	UploadLimitBytes int64
	DefaultGetLimit  int64
	MaxGetLimit      int64
}

func ProvideConfig() Config {
	return Config{
		SaltLength:       viper.GetInt("SALT_LENGTH"),
		PathUpload:       viper.GetString("PATH_UPLOAD"),
		UploadLimitBytes: viper.GetInt64("UPLOAD_LIMIT_BYTES"),
		DefaultGetLimit:  viper.GetInt64("DEFAULT_GET_LIMIT"),
		MaxGetLimit:      viper.GetInt64("MAX_GET_LIMIT"),
	}
}
