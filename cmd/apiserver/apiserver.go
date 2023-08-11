package main

/*
	The API server is the service that manages both the LittleBigPlanet API, and the new API for querying data

*/

import (
	"HugeSpaceship/pkg/api/game_api"
	"HugeSpaceship/pkg/api/web_api"
	"HugeSpaceship/pkg/common/config"
	"HugeSpaceship/pkg/common/db"
	"HugeSpaceship/pkg/common/logger"
	_ "embed"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"time"
)

// main is the entrypoint for the API server
func main() {
	startTime := time.Now()
	err := config.LoadConfig("apiserver")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	logger.LoggingInit("apiserver")

	_ = db.GetConnection()

	ctx := gin.New()
	ctx.Use(logger.LoggingMiddleware())

	api := ctx.Group("/api")

	// LittleBigPlanet compatible API
	if viper.GetBool("enable_gameserver") {
		gameAPI := api.Group("/LBP_XML")
		game_api.APIBootstrap(gameAPI)
	}
	// Web API
	if viper.GetBool("enable_api") {
		web_api.APIBootstrap(api)
	}

	log.Info().Float64("startSecs", time.Since(startTime).Seconds()).Msg("Time to start")
	err = ctx.Run("0.0.0.0:80")
	if err != nil {
		panic(err)
	}
}
