package config

import "github.com/cristalhq/aconfig"
import "github.com/cristalhq/aconfig/aconfigyaml"

type Config struct {
	Port     int `default:"8080" usage:"The listen port for the HTTP server"`
	Database struct {
		Host     string `required:"true"`
		Port     uint16 `required:"true"`
		Username string `required:"true"`
		Password string `required:"true"`
		Database string `required:"true" usage:"The Database to use"`
	} `usage:"Config for a postgresql database"`
	LBPApi struct {
		EnforceDigest      bool   `default:"false"`
		DigestKey          string ``
		AlternateDigestKey string
	}
	Log struct {
		Debug bool `default:"false"`
	}
}

func LoadConfig() (cfg *Config, err error) {
	cfg = &Config{}
	loader := aconfig.LoaderFor(cfg, aconfig.Config{
		EnvPrefix: "HS",
		Files:     []string{"/etc/hugespaceship/config.yml", "hugespaceship.yml"},
		FileDecoders: map[string]aconfig.FileDecoder{
			".yml": aconfigyaml.New(),
		},
	})

	loader.Flags()

	err = loader.Load()

	return cfg, err
}
