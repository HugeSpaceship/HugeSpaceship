package logger

import (
	"HugeSpaceship/internal/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

func LoggingInit(service string, config *config.Config) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	if config.Log.Debug {
		log.Logger = log.With().Caller().Logger().Output(zerolog.ConsoleWriter{Out: os.Stderr})
	} else {
		file, _ := os.OpenFile("./hs-"+service+".log", os.O_CREATE|os.O_WRONLY, 0644)
		log.Logger = log.Level(zerolog.InfoLevel).With().Str("service", service).Logger().Output(file)

	}
}
