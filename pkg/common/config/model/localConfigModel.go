package model

type LocalConfig struct {
	DBConfig struct {
		Host       string `yaml:"host"`
		Port       string `yaml:"port"`
		Username   string `yaml:"username"`
		Password   string `yaml:"password"`
		SSLEnabled bool   `yaml:"ssl_enabled,omitempty"`
	}
	ListenAddr string `yaml:"listen_addr"`
	ListenPort string `yaml:"listen_port"`
	PublicURL  string `yaml:"public_url"`
	SSLConfig  struct {
	} `yaml:"ssl_config,omitempty"`
}
