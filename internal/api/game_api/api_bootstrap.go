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
	"github.com/go-chi/chi/v5"
)

func ResourceBootstrap(group chi.Router, cfg *config.Config) {
	group.With(middlewares.TicketAuthMiddleware()).Get("/r/{hash}", resources.GetResourceHandler(cfg))
}

func APIBootstrap(cfg *config.Config) func(r chi.Router) {
	return func(r chi.Router) {
		r.Use(middlewares.PSPVersionMiddleware, middlewares.ServerHeaderMiddleware)

		r.Post("/login", auth.LoginHandler())
		r.With(middlewares.DigestMiddleware(cfg)).Get("/eula", controllers.EulaHandler())

		// LittleBigPlanet compatible API, required NpTicket auth
		authGameAPI := r.With(middlewares.TicketAuthMiddleware())
		authGameAPI.Get("/announce", controllers.AnnounceHandler())
		authGameAPI.Get("/network_settings.nws", settings.NetSettingsHandler())

		// LittleBigPlanet compatible API with digest calculation
		digestRequiredAPI := authGameAPI.With(middlewares.DigestMiddleware(cfg))
		digestRequiredAPI.Get("/user/{username}", users.UserGetHandler())
		digestRequiredAPI.Post("/match", match.MatchEndpoint())
		digestRequiredAPI.Post("/npdata", settings.NpDataEndpoint())
		digestRequiredAPI.Get("/notification", controllers.NotificationController()) // Stub
		digestRequiredAPI.Post("/goodbye", auth.LogoutHandler())
		digestRequiredAPI.Get("/news", controllers.NewsHandler())
		digestRequiredAPI.Get("/news/{id}", controllers.LBP2NewsHandler())
		digestRequiredAPI.Get("/stream", controllers.StreamHandler())

		digestRequiredAPI.Get("/slots", slots.GetSlotsHandler())
		digestRequiredAPI.Get("/slots/by", slots.GetSlotsByHandler())

		digestRequiredAPI.Get("/slots/lbp2luckydip", slots.LuckyDipHandler())
		digestRequiredAPI.Get("/slots/thumbs", slots.HighestRatedLevelsHandler())
		digestRequiredAPI.Get("/slots/lbp2cool", slots.HighestRatedLevelsHandler())
		digestRequiredAPI.Get("/slots/cool", slots.HighestRatedLevelsHandler())
		digestRequiredAPI.Get("/slots/highestRated", slots.HighestRatedLevelsHandler())
		digestRequiredAPI.Get("/s/user/{id}", slots.GetSlotHandler())

		digestRequiredAPI.Post("/scoreboard/{levelType}/{levelID}", slots.UploadScoreHandler())

		digestRequiredAPI.Post("/startPublish", slots.StartPublishHandler())
		digestRequiredAPI.Post("/publish", slots.PublishHandler())
		digestRequiredAPI.Post("/unpublish/{id}", slots.UnPublishHandler())
		digestRequiredAPI.Post("/updateUser", users.UpdateUserHandler())

		digestRequiredAPI.Post("/showModerated", moderation.ShowModeratedHandler())
		digestRequiredAPI.Post("/filter", moderation.FilterHandler())

		digestRequiredAPI.Post("/uploadPhoto", photos.UploadPhoto())
		digestRequiredAPI.Post("/photos/by", photos.GetPhotosBy())

		digestRequiredAPI.Post("/showNotUploaded", resources.ShowNotUploadedHandler())
		digestRequiredAPI.Post("/filterResources", resources.ShowNotUploadedHandler())
		digestRequiredAPI.Post("/upload/{hash}", resources.UploadResources())

		//Stubby, mc stub face
		digestRequiredAPI.Get("/promotions", controllers.StubEndpoint())
		digestRequiredAPI.Get("/user/{username}/playlists", controllers.StubEndpoint())
	}
}
