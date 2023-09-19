package config

import "github.com/cristalhq/aconfig"
import "github.com/cristalhq/aconfig/aconfigyaml"

// Config is the struct that contains all the global service config for the various components of the application
type Config struct {
	HTTPPort int `default:"8080" usage:"The listen port for the HTTP server"`
	Database struct {
		Host     string `required:"true"`
		Port     uint16 `required:"true"`
		Username string `required:"true"`
		Password string `required:"true"`
		Database string `required:"true" usage:"The Database to use"`
	} `usage:"Config for a postgresql database" env:"DB"`
	LBPApi struct {
		EnforceDigest      bool   `default:"false"`
		DigestKey          string ``
		AlternateDigestKey string
	}
	Log struct {
		Debug bool `default:"false"`
	}
}

// LoadConfig loads the configuration from various locations and returns a pointer to a Config struct
func LoadConfig(skipEnv bool) (cfg *Config, err error) {
	cfg = &Config{}
	loader := aconfig.LoaderFor(cfg, aconfig.Config{
		SkipFlags: true,
		SkipEnv:   skipEnv,
		EnvPrefix: "HS",
		Files:     []string{"/etc/hugespaceship/config.yml", "hugespaceship.yml"},
		FileDecoders: map[string]aconfig.FileDecoder{
			".yml": aconfigyaml.New(),
		},
	})

	err = loader.Load()
	return cfg, err
}
