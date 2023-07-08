package logger

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"os"
)

func LoggingInit(service string) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	if viper.GetBool("debug") {
		log.Logger = log.With().Caller().Logger().Output(zerolog.ConsoleWriter{Out: os.Stderr})
	} else {
		file, _ := os.OpenFile("./hs-"+service+".log", os.O_CREATE|os.O_WRONLY, 0644)
		log.Logger = log.Level(zerolog.InfoLevel).With().Str("service", service).Logger().Output(file)

	}
}
