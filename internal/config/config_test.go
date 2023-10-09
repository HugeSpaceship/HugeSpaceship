package config

import (
	"fmt"
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {

	t.Run("Environment config test", testEnvConfig)
	t.Run("Environment config test with invalid data", testInvalidEnvConfig)

	err := os.Chdir("../../test") // So we can load the config files
	if err != nil {
		t.Fatal("Can't chdir to config dir")
	}

	t.Run("Config file test", testConfigFile)
}

func testEnvConfig(t *testing.T) {
	t.Setenv("HS_HTTP_PORT", "8080")

	t.Setenv("HS_DB_PORT", "5432")
	t.Setenv("HS_DB_HOST", "testHost")
	t.Setenv("HS_DB_USERNAME", "testUsername")
	t.Setenv("HS_DB_PASSWORD", "testPassword")
	t.Setenv("HS_DB_DATABASE", "testDB")

	cfg, err := LoadConfig(false)
	if err != nil {
		t.Fatal("Failed to load config: ", err)
	}

	if cfg.HTTPPort != 8080 {
		t.Error("Unexpected value, expected 8080 ", cfg.HTTPPort)
	}

	err = os.Unsetenv("HS_HTTP_PORT")
	if err != nil {
		t.Fatal("couldn't unset environment variable", err)
	}
}

func testInvalidEnvConfig(t *testing.T) {
	t.Setenv("HS_HTTP_PORT", "a lizard")

	fmt.Println(os.Getenv("HS_HTTP_PORT"))

	cfg, err := LoadConfig(false)
	if err == nil {
		t.Fatal("The config loaded in a situation where it shouldn't")
	}

	if cfg.HTTPPort != 8080 {
		t.Error("Unexpected value, it should be the default", cfg.HTTPPort)
	}
}

func testConfigFile(t *testing.T) {
	cfg, err := LoadConfig(true)
	if err != nil {
		t.Fatal("Failed to load config", err)
	}

	if cfg == nil {
		t.Fatal("Config pointer is nil")
	}

	if cfg.HTTPPort != 8181 {
		t.Error("Unexpected value, expected 8181", cfg.HTTPPort)
	}

	if cfg.API.EnforceDigest != true {
		t.Error("Unexpected value, expected true", cfg.API.EnforceDigest)
	}
	if cfg.API.DigestKey != "test_digestKey" {
		t.Error("Unexpected value", cfg.API.DigestKey)
	}
	if cfg.API.AlternateDigestKey != "test_altDigestKey" {
		t.Error("Unexpected value", cfg.API.AlternateDigestKey)
	}
	if cfg.Database.Database != "testDB" {
		t.Error("Unexpected value")
	}
	if cfg.Database.Host != "testHost" {
		t.Error("Unexpected value, expected testHost", cfg.Database.Host)
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
