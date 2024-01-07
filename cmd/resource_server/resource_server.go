package main

import (
	"HugeSpaceship/internal/api/game_api"
	"HugeSpaceship/internal/config"
	"HugeSpaceship/pkg/db"
	"HugeSpaceship/pkg/db/migration"
	"HugeSpaceship/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
)

func main() {
	cfg, err := config.LoadConfig(false)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	logger.LoggingInit("api_server", cfg)

	pool := db.Open(cfg)            // Open a connection to the DB
	err = migration.MigrateDB(pool) // Migrate the DB to the latest schema
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to migrate database")
	}

	ctx := gin.New()
	ctx.Use(logger.LoggingMiddleware())

	// everything starts at /api
	api := ctx.Group("/api")
	game_api.ResourceBootstrap(api, cfg)

	err = ctx.Run("0.0.0.0:80")
	if err != nil {
		panic(err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
}
