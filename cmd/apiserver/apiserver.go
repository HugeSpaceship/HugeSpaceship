package main

/*
	The API server is the service that manages both the LittleBigPlanet API, and the new API for querying data

*/

import (
	"HugeSpaceship/pkg/api/game_api"
	"HugeSpaceship/pkg/api/web_api"
	"HugeSpaceship/pkg/common/config"
	"HugeSpaceship/pkg/common/db"
	"HugeSpaceship/pkg/common/db/migration"
	"HugeSpaceship/pkg/common/logger"
	_ "embed"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"time"
)

// main is the entrypoint for the API server
func main() {
	// So we can time how long it took to start the server
	startTime := time.Now()

	cfg, err := config.LoadConfig(false)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	logger.LoggingInit("apiserver", cfg)

	pool := db.Open(cfg) // Open a connection to the DB
	err = migration.MigrateDB(pool)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to migrate database")
	}

	ctx := gin.New()
	ctx.Use(logger.LoggingMiddleware())

	api := ctx.Group("/api")

	// LittleBigPlanet compatible API
	game_api.APIBootstrap(api, cfg)

	// Web API
	web_api.APIBootstrap(api)

	log.Info().Float64("startSecs", time.Since(startTime).Seconds()).Msg("Time to start")
	err = ctx.Run("0.0.0.0:80")
	if err != nil {
		panic(err)
	}
}
