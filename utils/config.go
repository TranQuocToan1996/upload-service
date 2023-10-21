package utils

import (
	"errors"
	"fmt"
	"os"

	"github.com/labstack/gommon/log"
	"github.com/ory/viper"
)

func InitConfig(path string) {
	viper.AutomaticEnv()
	viper.SetConfigFile(path)
	err := viper.ReadInConfig()
	var pathErr *os.PathError
	if errors.As(err, &pathErr) {
		log.Warnf("no config file '%s' not found. Using default values", path)
	} else if err != nil { // Handle other errors that occurred while reading the config file
		panic(fmt.Errorf("fatal error while reading the config file: %w", err))
	}
}
