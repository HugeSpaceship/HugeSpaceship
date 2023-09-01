package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	err := os.Chdir("../../../test") // So we can load the config files
	if err != nil {
		t.Fatal("Can't chdir to config dir")
	}

	t.Run("Config file test", TestConfigFile)
	t.Run("Environment config test", TestEnvConfig)
	t.Run("Environment config test with invalid data", TestEnvConfig)
}

func TestEnvConfig(t *testing.T) {
	err := os.Setenv("HS_HTTP_PORT", "8080")
	if err != nil {
		t.Fatal("couldn't set environment variable")
	}
	cfg, err := LoadConfig()
	if err != nil {
		t.Fatal("Failed to load config")
	}

	if cfg.HTTPPort != 8080 {
		t.Error("Unexpected value")
	}
}

func TestInvalidEnvConfig(t *testing.T) {
	err := os.Setenv("HS_HTTP_PORT", "a lizard")
	if err != nil {
		t.Fatal("couldn't set environment variable")
	}
	cfg, err := LoadConfig()
	if err == nil {
		t.Fatal("The config loaded in a situation where it shouldn't")
	}

	if cfg.HTTPPort != 0 {
		t.Error("Unexpected value")
	}
}

func TestConfigFile(t *testing.T) {
	cfg, err := LoadConfig()
	if err != nil {
		t.Fatal("Failed to load config")
	}

	if cfg == nil {
		t.Fatal("Config pointer is nil")
	}

	if cfg.HTTPPort != 8080 {
		t.Error("Unexpected value")
	}

	if cfg.LBPApi.EnforceDigest != true {
		t.Error("Unexpected value")
	}
	if cfg.LBPApi.DigestKey != "test_digestKey" {
		t.Error("Unexpected value")
	}
	if cfg.LBPApi.AlternateDigestKey != "test_altDigestKey" {
		t.Error("Unexpected value")
	}
	if cfg.Database.Database != "testDB" {
		t.Error("Unexpected value")
	}
	if cfg.Database.Host != "testHost" {
		t.Error("Unexpected value")
	}
	if cfg.Database.Port != 8425 {
		t.Error("Unexpected value")
	}
	if cfg.Database.Username != "testUsername" {
		t.Error("Unexpected value")
	}

	if cfg.Database.Password != "testPassword" {
		t.Error("Unexpected value")
	}
}
