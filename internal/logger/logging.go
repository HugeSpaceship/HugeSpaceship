package logger

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"os"
)

func LoggingInit(service string, v *viper.Viper) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	if v.GetBool("log.debug") {
		log.Logger = log.With().Caller().Logger().Output(zerolog.ConsoleWriter{Out: os.Stderr})
	} else {
		v.SetDefault("log.level", "info")
		level, err := zerolog.ParseLevel(v.GetString("log.level"))
		if err != nil {
			panic(err)
		}
		logger := log.Level(level).With()
		if service != "" {
			logger = logger.Str("service", service)
		}

		if v.GetBool("log.log-to-file") {
			file, _ := os.OpenFile("./hs-"+service+".log", os.O_CREATE|os.O_WRONLY, 0644)
			log.Logger = logger.Logger().Output(file)
		} else {
			log.Logger = logger.Logger()
		}
	}
}
