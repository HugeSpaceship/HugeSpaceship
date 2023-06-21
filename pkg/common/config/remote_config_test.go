package config

import (
	"HugeSpaceship/pkg/common/config/model/remote"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"testing"
)

func TestGetLBPAPIConfig(t *testing.T) {
	viper.SetEnvPrefix("hs")
	viper.AutomaticEnv()
	GetLBPAPIConfig()
}

func TestSaveLBPAPIConfig(t *testing.T) {
	viper.SetEnvPrefix("hs")
	viper.AutomaticEnv()
	config := remote.LBPAPIConfig{
		PrimaryDigest:   "test1",
		AlternateDigest: "test2",
		AutoCreateUsers: "true",
	}
	err := SaveLBPAPIConfig(config)
	log.Error().Err(err).Msg("idk dude?")
}
