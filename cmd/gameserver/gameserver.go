package main

import (
	"HugeSpaceship/pkg/game_server/controllers"
	"HugeSpaceship/pkg/game_server/controllers/auth"
	"HugeSpaceship/pkg/logger"
	"github.com/gin-gonic/gin"
)

func main() {
	ctx := gin.New()
	ctx.Use(gin.Logger(), logger.LoggingMiddleware())
	api := ctx.Group("/LITTLEBIGPLANETPS3_XML")
	api.POST("/login", auth.LoginHandler())
	api.GET("/eula", controllers.EulaHandler())
	api.GET("/announce", controllers.AnnounceHandler())

	err := ctx.Run(":80")
	if err != nil {
		panic(err)
	}
}
