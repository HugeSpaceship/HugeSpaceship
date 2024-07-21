package config

import (
	"fmt"
	"github.com/spf13/viper"
	"strings"
)

type ResourceBackendConfig struct {
	Name     string                 `yaml:"name"`
	Type     string                 `yaml:"type"`
	Priority uint                   `yaml:"priority"`
	Config   map[string]interface{} `yaml:"config"`
}

// Config is the struct that contains all the global service config for the various components of the application
type Config struct {
	HTTPPort int `default:"8080" usage:"The listen port for the HTTP server" env:"HTTP_PORT"`
	Database struct {
		Host     string `required:"true"`
		Port     uint16 `required:"true"`
		Username string `required:"true"`
		Password string `required:"true"`
		Database string `required:"true" usage:"The Database to use"`
	} `usage:"Config for a postgresql database" env:"DB"`
	API struct {
		EnforceDigest      bool `default:"false"`
		DigestKey          string
		AlternateDigestKey string
	}
	ResourceServer struct {
		Enabled        bool   `default:"true"`
		CacheResources bool   `default:"false"`
		CacheLocation  string `default:"./r"`
		Backends       map[string]*ResourceBackendConfig
	}
	Website struct {
		Enabled              bool   `default:"true"`
		UseEmbeddedResources bool   `default:"true"`
		WebRoot              string `env:"WEBROOT"`
		ThemePath            string
		DefaultTheme         string
		AllowUserThemes      bool
	}
	Log struct {
		Debug       bool   `default:"false"` // Debug controls the web server debug
		FileLogging bool   `default:"true"`
		Level       string `default:"info"`
	}
}

// LoadConfig loads the configuration from various locations and returns a pointer to a Config struct
func LoadConfig(skipEnv bool) (v *viper.Viper) {

	v = viper.New()

	v.SetConfigName("hugespaceship")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("/etc/hugespaceship")
	err := v.ReadInConfig() // Find and read the config file
	if err != nil {         // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	v.SetEnvPrefix("hs")

	if !skipEnv {
		v.SetEnvKeyReplacer(strings.NewReplacer("_", "-"))
		v.AutomaticEnv()
	}

	//cfg = &Config{}
	//loader := aconfig.LoaderFor(cfg, aconfig.Config{
	//	SkipFlags: true,
	//	SkipEnv:   skipEnv,
	//	EnvPrefix: "HS",
	//	Files:     []string{"/etc/hugespaceship/config.yml", "hugespaceship.yml"},
	//	FileDecoders: map[string]aconfig.FileDecoder{
	//		".yml": aconfigyaml.New(),
	//	},
	//})
	//
	//err = loader.Load()
	return v
}
