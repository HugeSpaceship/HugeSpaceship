package logger

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"os"
)

func LoggingInit(service string) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	log.Logger = log.With().Str("service", service).Caller().Logger()

	if viper.GetBool("debug") {
		log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
}
