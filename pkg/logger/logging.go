package logger

import (
	"github.com/HugeSpaceship/HugeSpaceship/internal/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

func LoggingInit(service string, config *config.Config) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	if config.Log.Debug {
		log.Logger = log.With().Caller().Logger().Output(zerolog.ConsoleWriter{Out: os.Stderr})
	} else {
		level, err := zerolog.ParseLevel(config.Log.Level)
		if err != nil {
			panic(err)
		}
		logger := log.Level(level).With()
		if service != "" {
			logger = logger.Str("service", service)
		}

		if config.Log.FileLogging {
			file, _ := os.OpenFile("./hs-"+service+".log", os.O_CREATE|os.O_WRONLY, 0644)
			log.Logger = logger.Logger().Output(file)
		} else {
			log.Logger = logger.Logger()
		}
	}
}
