package game_api

import (
	"HugeSpaceship/internal/api/game_api/controllers"
	"HugeSpaceship/internal/api/game_api/controllers/auth"
	"HugeSpaceship/internal/api/game_api/controllers/match"
	"HugeSpaceship/internal/api/game_api/controllers/moderation"
	"HugeSpaceship/internal/api/game_api/controllers/photos"
	"HugeSpaceship/internal/api/game_api/controllers/resources"
	"HugeSpaceship/internal/api/game_api/controllers/settings"
	"HugeSpaceship/internal/api/game_api/controllers/slots"
	"HugeSpaceship/internal/api/game_api/controllers/users"
	"HugeSpaceship/internal/api/game_api/middlewares"
	"HugeSpaceship/internal/config"
	"github.com/gin-gonic/gin"
)

func ResourceBootstrap(group *gin.RouterGroup, cfg *config.Config) {
	group.GET("/r/:hash", resources.GetResourceHandler(cfg), middlewares.TicketAuthMiddleware())
}

func APIBootstrap(gameAPI *gin.RouterGroup, cfg *config.Config) {
	gameAPI.Use(middlewares.PSPVersionMiddleware, middlewares.ServerHeaderMiddleware)

	gameAPI.POST("/login", auth.LoginHandler())
	gameAPI.GET("/eula", middlewares.DigestMiddleware(cfg), controllers.EulaHandler())

	// LittleBigPlanet compatible API, required NpTicket auth
	authGameAPI := gameAPI.Group("", middlewares.TicketAuthMiddleware())
	authGameAPI.GET("/announce", controllers.AnnounceHandler())
	authGameAPI.GET("/network_settings.nws", settings.NetSettingsHandler())

	// LittleBigPlanet compatible API with digest calculation
	digestRequiredAPI := authGameAPI.Group("", middlewares.DigestMiddleware(cfg))
	digestRequiredAPI.GET("/user/:username", users.UserGetHandler())
	digestRequiredAPI.POST("/match", match.MatchEndpoint())
	digestRequiredAPI.POST("/npdata", settings.NpDataEndpoint())
	digestRequiredAPI.GET("/notification", controllers.NotificationController()) // Stub
	digestRequiredAPI.POST("/goodbye", auth.LogoutHandler())
	digestRequiredAPI.GET("/news", controllers.NewsHandler())
	digestRequiredAPI.GET("/news/:id", controllers.LBP2NewsHandler())
	digestRequiredAPI.GET("/stream", controllers.StreamHandler())

	digestRequiredAPI.GET("/slots", slots.GetSlotsHandler())
	digestRequiredAPI.GET("/slots/by", slots.GetSlotsByHandler())

	digestRequiredAPI.GET("/slots/lbp2luckydip", slots.LuckyDipHandler())
	digestRequiredAPI.GET("/slots/thumbs", slots.HighestRatedLevelsHandler())
	digestRequiredAPI.GET("/slots/lbp2cool", slots.HighestRatedLevelsHandler())
	digestRequiredAPI.GET("/slots/cool", slots.HighestRatedLevelsHandler())
	digestRequiredAPI.GET("/slots/highestRated", slots.HighestRatedLevelsHandler())
	digestRequiredAPI.GET("/s/user/:id", slots.GetSlotHandler())

	digestRequiredAPI.POST("/scoreboard/:levelType/:levelID", slots.UploadScoreHandler())

	digestRequiredAPI.POST("/startPublish", slots.StartPublishHandler())
	digestRequiredAPI.POST("/publish", slots.PublishHandler())
	digestRequiredAPI.POST("/unpublish/:id", slots.UnPublishHandler())
	digestRequiredAPI.POST("/updateUser", users.UpdateUserHandler())

	digestRequiredAPI.POST("/showModerated", moderation.ShowModeratedHandler())
	digestRequiredAPI.POST("/filter", moderation.FilterHandler())

	digestRequiredAPI.POST("/uploadPhoto", photos.UploadPhoto())
	digestRequiredAPI.POST("/photos/by", photos.GetPhotosBy())

	digestRequiredAPI.POST("/showNotUploaded", resources.ShowNotUploadedHandler())
	digestRequiredAPI.POST("/filterResources", resources.ShowNotUploadedHandler())
	digestRequiredAPI.POST("/upload/:hash", resources.UploadResources())

	//Stubby, mc stub face
	digestRequiredAPI.GET("/promotions", controllers.StubEndpoint())
	digestRequiredAPI.GET("/user/:username/playlists", controllers.StubEndpoint())
}
