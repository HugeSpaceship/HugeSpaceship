package config

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

// LoadConfig reads the config for the specified service
func LoadConfig(service string) error {
	viper.SetConfigName(service)
	viper.SetConfigType("toml")
	viper.AddConfigPath("/etc/hugespaceship/")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		log.Info().Msg("no config file found, if you are loading config from the environment this is fine")
	} else {
		return err
	}
	viper.SetEnvPrefix("hs_" + service)
	viper.AutomaticEnv()
	viper.SetDefault("db_host", "localhost")
	return nil
}
