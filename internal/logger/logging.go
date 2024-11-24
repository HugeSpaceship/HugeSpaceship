package logger

import (
	"errors"
	"github.com/HugeSpaceship/HugeSpaceship/internal/config"
	"io"
	"log/slog"
	"os"
)

var defaultLevel = slog.LevelInfo

func LoggingInit(cfg *config.Config) error {
	level := defaultLevel
	switch cfg.Log.Level {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		return errors.New("invalid log level")
	}

	slog.SetLogLoggerLevel(level)

	var logger *slog.Logger

	var writers = []io.Writer{os.Stdout}

	if cfg.Log.LogFile != "" {
		logFile, err := os.OpenFile(cfg.Log.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}

		writers = append(writers, logFile)
	}

	if cfg.Log.JSONLogging {
		jsonHandler := slog.NewJSONHandler(io.MultiWriter(writers...), &slog.HandlerOptions{
			AddSource:   cfg.Log.DebugInfo,
			Level:       level,
			ReplaceAttr: nil,
		})

		logger = slog.New(jsonHandler)

	} else {
		textHandler := slog.NewTextHandler(io.MultiWriter(writers...), &slog.HandlerOptions{
			AddSource:   cfg.Log.DebugInfo,
			Level:       level,
			ReplaceAttr: nil,
		})

		logger = slog.New(textHandler)
	}

	slog.SetDefault(logger)
	return nil
}
