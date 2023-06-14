package config

import (
	"HugeSpaceship/pkg/common/config/model/remote"
	"github.com/rs/zerolog/log"
	"testing"
)

func TestGetLBPAPIConfig(t *testing.T) {

	GetLBPAPIConfig()
}

func TestSaveLBPAPIConfig(t *testing.T) {
	config := remote.LBPAPIConfig{
		PrimaryDigest:   "test1",
		AlternateDigest: "test2",
		AutoCreateUsers: "true",
	}
	err := SaveLBPAPIConfig(config)
	log.Error().Err(err).Msg("idk dude?")
}
