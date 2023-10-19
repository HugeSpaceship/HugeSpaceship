package main

/*
	The API server is the service that manages both the LittleBigPlanet API, and the new API for querying data

*/

import (
	"HugeSpaceship/internal/api/game_api"
	"HugeSpaceship/internal/config"
	"HugeSpaceship/internal/website"
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

	logger.LoggingInit("apiserver", cfg)

	pool := db.Open(cfg)            // Open a connection to the DB
	err = migration.MigrateDB(pool) // Migrate the DB to the latest schema
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to migrate database")
	}

	// Init the web framework
	ctx := gin.New()
	ctx.Use(logger.LoggingMiddleware())

	// everything starts at /api
	api := ctx.Group("/LITTLEGIBPLANETPS3_XML")
	api2 := ctx.Group("/LITTLEBIGPLANETPS3_XML")
	api3 := ctx.Group("/api/LBP_XML")
	// LittleBigPlanet compatible API
	game_api.APIBootstrap(api, cfg)

	// LittleBigPlanet compatible API, but with brainfart this time
	game_api.APIBootstrap(api2, cfg)
	game_api.APIBootstrap(api3, cfg)
	// Resource server
	if cfg.ResourceServer.Enabled {
		game_api.ResourceBootstrap(api, cfg)
		game_api.ResourceBootstrap(api2, cfg)
	}

	if cfg.Website.Enabled {
		website.Bootstrap(ctx, cfg)
	}

	err = ctx.Run("0.0.0.0:" + strconv.Itoa(cfg.HTTPPort))
	if err != nil {
		panic(err)
	}
}
