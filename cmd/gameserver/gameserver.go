package main

import (
	"HugeSpaceship/pkg/game_server/controllers"
	"HugeSpaceship/pkg/game_server/controllers/auth"
	"HugeSpaceship/pkg/logger"
	_ "embed"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

//go:generate sh -c "printf %s $(git rev-parse --short HEAD) > VERSION.txt"
//go:embed VERSION.txt
var commitHash string

func main() {

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Info().Str("commitHash", commitHash)

	ctx := gin.New()
	ctx.Use(logger.LoggingMiddleware())
	api := ctx.Group("/LITTLEBIGPLANETPS3_XML")
	api.POST("/login", auth.LoginHandler())
	api.GET("/eula", controllers.EulaHandler())
	api.GET("/announce", controllers.AnnounceHandler())

	err := ctx.Run(":80")
	if err != nil {
		panic(err)
	}
}
