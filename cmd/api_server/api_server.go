package main

/*
	The API server is the service that manages both the LittleBigPlanet API, and the new API for querying data

*/

import (
	"HugeSpaceship/internal/api/game_api"
	"HugeSpaceship/internal/config"
	"HugeSpaceship/pkg/db"
	"HugeSpaceship/pkg/db/migration"
	"HugeSpaceship/pkg/logger"
	_ "embed"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"strconv"
)

// main is the entrypoint for the API server
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

	// Init the web framework
	ctx := gin.New()
	ctx.Use(logger.LoggingMiddleware())

	// everything starts at /api
	api := ctx.Group("/api/LBP_XML")
	// LittleBigPlanet compatible API
	game_api.APIBootstrap(api, cfg)

	err = ctx.Run("0.0.0.0:" + strconv.Itoa(cfg.HTTPPort))
	if err != nil {
		panic(err)
	}
}
