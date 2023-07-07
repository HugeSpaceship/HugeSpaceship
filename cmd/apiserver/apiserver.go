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
)

// main is the entrypoint for the API server
func main() {
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
	gameAPI := api.Group("/LBP_XML")
	game_api.APIBootstrap(gameAPI)

	// Web API, optional
	if viper.GetBool("enable_api") {
		web_api.APIBootstrap(api)
	}

	err = ctx.Run(":80")
	if err != nil {
		panic(err)
	}
}
