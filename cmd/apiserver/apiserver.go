package main

/*
	The API server is the service that manages both the LittleBigPlanet API, and the new API for querying data

*/

import (
	"HugeSpaceship/pkg/api/game_api/controllers"
	"HugeSpaceship/pkg/api/game_api/controllers/auth"
	"HugeSpaceship/pkg/api/game_api/controllers/match"
	"HugeSpaceship/pkg/api/game_api/controllers/resources"
	"HugeSpaceship/pkg/api/game_api/controllers/settings"
	"HugeSpaceship/pkg/api/game_api/controllers/slots"
	"HugeSpaceship/pkg/api/game_api/controllers/users"
	"HugeSpaceship/pkg/api/game_api/middlewares"
	"HugeSpaceship/pkg/common/db"
	"HugeSpaceship/pkg/common/logger"
	_ "embed"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

//go:generate sh -c "printf %s $(git rev-parse --short HEAD) > VERSION.txt"
//go:embed VERSION.txt

var commitHash string

// main is the entrypoint for the API server
func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Info().Str("commitHash", commitHash).Msg("Starting HugeSpaceship API Server")
	gin.Logger()
	_ = db.GetConnection()

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

	// LittleBigPlanet compatible API with digest calculation
	digestRequiredAPI := authGameAPI.Group("", middlewares.DigestMiddleware())
	digestRequiredAPI.GET("/user/:username", users.UserGetHandler())
	digestRequiredAPI.POST("/match", match.MatchEndpoint())
	digestRequiredAPI.POST("/npdata", settings.NpDataEndpoint())
	digestRequiredAPI.GET("/notification", controllers.NotificationController()) // Stub
	digestRequiredAPI.POST("/goodbye", auth.LogoutHandler())
	digestRequiredAPI.GET("/news", controllers.NewsHandler())
	digestRequiredAPI.GET("/news/:id", controllers.LBP2NewsHandler())
	digestRequiredAPI.GET("/stream", controllers.StreamHandler())
	digestRequiredAPI.POST("/startPublish", slots.StartPublishHandler())
	digestRequiredAPI.POST("/upload/:hash", resources.UploadResources())
	digestRequiredAPI.GET("/slots/by", slots.GetSlotsByHandler())
	// Web API
	webAPI := api.Group("/v1")
	webAPI.GET("/users")

	err := ctx.Run(":80")
	if err != nil {
		panic(err)
	}
}
