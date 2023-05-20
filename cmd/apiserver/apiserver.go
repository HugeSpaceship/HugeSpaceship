package main

/*
	The API server is the service that manages both the LittleBigPlanet API, and the new API for querying data

*/

import (
	"HugeSpaceship/pkg/api/game_api/controllers"
	"HugeSpaceship/pkg/api/game_api/controllers/auth"
	"HugeSpaceship/pkg/api/game_api/middlewares"
	"HugeSpaceship/pkg/common/logger"
	_ "embed"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

//go:generate sh -c "printf %s $(git rev-parse --short HEAD) > VERSION.txt"
//go:embed VERSION.txt
var commitHash string

// main is the entrypoint for the API server
func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Info().Str("commitHash", commitHash)

	ctx := gin.New()
	ctx.Use(logger.LoggingMiddleware())

	api := ctx.Group("/api")

	// LittleBigPlanet compatible API
	gameAPI := api.Group("/LBP_XML")
	gameAPI.POST("/login", auth.LoginHandler())
	gameAPI.GET("/eula", controllers.EulaHandler())

	// LittleBigPlanet compatible API, required NpTicket auth
	authGameAPI := gameAPI.Group("", middlewares.TicketAuthMiddleware())
	authGameAPI.GET("/announce", controllers.AnnounceHandler())

	// Web API
	webAPI := api.Group("/v1")
	webAPI.GET("/users")

	err := ctx.Run(":80")
	if err != nil {
		panic(err)
	}
}
