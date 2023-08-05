package game_api

import (
	"HugeSpaceship/pkg/api/game_api/controllers"
	"HugeSpaceship/pkg/api/game_api/controllers/auth"
	"HugeSpaceship/pkg/api/game_api/controllers/match"
	"HugeSpaceship/pkg/api/game_api/controllers/resources"
	"HugeSpaceship/pkg/api/game_api/controllers/settings"
	"HugeSpaceship/pkg/api/game_api/controllers/slots"
	"HugeSpaceship/pkg/api/game_api/controllers/users"
	"HugeSpaceship/pkg/api/game_api/middlewares"
	"github.com/gin-gonic/gin"
)

func APIBootstrap(group *gin.RouterGroup) {
	group.POST("/login", auth.LoginHandler())
	group.GET("/eula", controllers.EulaHandler())

	// LittleBigPlanet compatible API, required NpTicket auth
	authGameAPI := group.Group("", middlewares.TicketAuthMiddleware())
	authGameAPI.GET("/announce", controllers.AnnounceHandler())
	authGameAPI.GET("/r/:hash", resources.GetResourceHandler())

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
	digestRequiredAPI.GET("/s/user/:id", slots.GetSlotHandler())
	digestRequiredAPI.POST("/publish", slots.PublishHandler())
	digestRequiredAPI.POST("/updateUser", users.UpdateUserHandler())
}
