package config

import (
	"HugeSpaceship/pkg/common/config/model"
	"gopkg.in/yaml.v3"
	"os"
	"sync"
)

var loadedConfig = model.LocalConfig{}
var runOnce = sync.Once{}

func GetConfig() model.LocalConfig {
	runOnce.Do(loadConfig)
	return loadedConfig
}

// Local config is for really basic things like the db config, remote config is where our distributed hot-reloadable settings should go
func loadConfig() {
	data, err := os.ReadFile("hugespaceship.yml")
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(data, &loadedConfig)
	if err != nil {
		panic(err)
	}
}
