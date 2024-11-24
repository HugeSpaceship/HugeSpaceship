package main

import (
	"github.com/HugeSpaceship/HugeSpaceship/internal/config"
	"github.com/HugeSpaceship/HugeSpaceship/internal/db"
	"github.com/HugeSpaceship/HugeSpaceship/internal/logger"
	"log/slog"
)

func main() {
	cfg, err := config.LoadConfig(false)
	if err != nil {
		slog.Error("Failed to load config", "error", err)
		panic(err)
	}

	err = logger.LoggingInit(cfg)
	if err != nil {
		slog.Error("Failed to initialize logger", "error", err)
	}

	_ = db.Open(cfg)
}
