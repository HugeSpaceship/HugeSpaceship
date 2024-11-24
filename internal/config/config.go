package config

import (
	"github.com/cristalhq/aconfig"
	"github.com/cristalhq/aconfig/aconfigyaml"
	"os"
)

// Config is the struct that contains all the global service config for the various components of the application
type Config struct {
	ListenAddr string `default:"0.0.0.0:10060" usage:"The listen address for the HTTP server" env:"LISTEN_ADDR"`
	Database   struct {
		Host     string `required:"true"`
		Port     uint16 `required:"true"`
		Username string `required:"true"`
		Password string `required:"true"`
		Database string `required:"true" usage:"The Database to use"`
	} `usage:"Config for a postgresql database" env:"DB"`
	GameAPI struct {
		EnforceDigest      bool   `default:"false"`
		DigestKey          string `env:"DIGEST_KEY"`
		AlternateDigestKey string `env:"ALTERNATE_DIGEST_KEY"`
	} `usage:"Config for the game api" env:"GAME_API"`
	ResourceServer struct {
		Backend string

		// S3/Object store connection info
		Endpoint        string
		BucketName      string
		Region          string
		AccessKeyID     string
		AccessKeySecret string

		// File store info
		ResourcePath string

		CacheResources bool   `default:"false"`
		CacheLocation  string `default:"./r"`
	}
	Website struct {
		Enabled              bool   `default:"true"`
		UseEmbeddedResources bool   `default:"true"`
		WebRoot              string `env:"WEBROOT"`
	}
	Log struct {
		LogFile     string `default:""`
		Level       string `default:"info"`
		JSONLogging bool   `default:"false"`
		DebugInfo   bool   `default:"false"`
	}
}

var CfgPriority = []string{
	"/usr/share/hugespaceship/config.yaml",
	"/usr/local/share/hugespaceship/config.yaml",
	"/etc/hugespaceship/config.yaml",
	"./config.yaml",
}

// LoadConfig loads the configuration from various locations and returns a pointer to a Config struct
func LoadConfig(skipEnv bool) (*Config, error) {

	cfg := Config{}

	loader := aconfig.LoaderFor(&cfg, aconfig.Config{
		SkipEnv:   skipEnv,
		EnvPrefix: "HS",

		Files: CfgPriority,
		FileDecoders: map[string]aconfig.FileDecoder{
			".yaml": aconfigyaml.New(),
		},
	})

	flag := loader.Flags()

	err := flag.Parse(os.Args[1:])
	if err != nil {
		return nil, err
	}

	err = loader.Load()

	return &cfg, err
}
