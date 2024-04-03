package game_api

import (
	"HugeSpaceship/internal/config"
	"HugeSpaceship/internal/http/api/game_api/controllers"
	"HugeSpaceship/internal/http/api/game_api/controllers/auth"
	"HugeSpaceship/internal/http/api/game_api/controllers/match"
	"HugeSpaceship/internal/http/api/game_api/controllers/moderation"
	"HugeSpaceship/internal/http/api/game_api/controllers/photos"
	"HugeSpaceship/internal/http/api/game_api/controllers/resources"
	"HugeSpaceship/internal/http/api/game_api/controllers/settings"
	"HugeSpaceship/internal/http/api/game_api/controllers/slots"
	"HugeSpaceship/internal/http/api/game_api/controllers/users"
	"HugeSpaceship/internal/http/api/game_api/middlewares"
	"HugeSpaceship/internal/http/middleware"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

func ResourceBootstrap(group chi.Router, cfg *config.Config) {
	group.With(middlewares.TicketAuthMiddleware).Get("/r/{hash}", resources.GetResourceHandler(cfg))
}

func APIBootstrap(r chi.Router, cfg *config.Config) {
	r.Use(middlewares.PSPVersionMiddleware, middlewares.ServerHeaderMiddleware)

	r.Use(chiMiddleware.RequestID)
	r.Use(chiMiddleware.Recoverer)
	r.Use(chiMiddleware.NoCache)

	r.With(middleware.DBCtxMiddleware).Post("/login", auth.LoginHandler())
	r.Get("/eula", controllers.EulaHandler())

	// LittleBigPlanet compatible API, required NpTicket auth
	authGameAPI := r.With(middleware.DBCtxMiddleware, middlewares.TicketAuthMiddleware)
	authGameAPI.With(middleware.DBCtxMiddleware).Get("/announce", controllers.AnnounceHandler())
	authGameAPI.Get("/network_settings.nws", settings.NetSettingsHandler())

	// LittleBigPlanet compatible API with digest calculation
	digestRequiredAPI := authGameAPI.With(middlewares.DigestMiddleware(cfg))
	digestRequiredAPI.Post("/match", match.MatchEndpoint())
	digestRequiredAPI.Post("/upload/{hash}", resources.UploadResources())

	// API with automatic content type
	xmlAPI := digestRequiredAPI.With(render.SetContentType(render.ContentTypeXML))
	xmlAPI.Get("/user/{username}", users.UserGetHandler())

	xmlAPI.Post("/npdata", settings.NpDataEndpoint())
	xmlAPI.Get("/notification", controllers.NotificationController()) // Stub
	xmlAPI.Post("/goodbye", auth.LogoutHandler())
	xmlAPI.Get("/news", controllers.NewsHandler())
	xmlAPI.Get("/news/{id}", controllers.LBP2NewsHandler())
	xmlAPI.Get("/stream", controllers.StreamHandler())

	xmlAPI.Route("/slots", func(r chi.Router) {
		xmlAPI.Get("/slots", slots.GetSlotsHandler())
		xmlAPI.Get("/slots/by", slots.GetSlotsByHandler())

		xmlAPI.Get("/slots/lbp2luckydip", slots.LuckyDipHandler())
		xmlAPI.Get("/slots/thumbs", slots.HighestRatedLevelsHandler())
		xmlAPI.Get("/slots/lbp2cool", slots.HighestRatedLevelsHandler())
		xmlAPI.Get("/slots/cool", slots.HighestRatedLevelsHandler())
		xmlAPI.Get("/slots/highestRated", slots.HighestRatedLevelsHandler())
	})

	xmlAPI.Get("/s/user/{id}", slots.GetSlotHandler())

	xmlAPI.Post("/scoreboard/{levelType}/{levelID}", slots.UploadScoreHandler())

	xmlAPI.Post("/startPublish", slots.StartPublishHandler())
	xmlAPI.Post("/publish", slots.PublishHandler())
	xmlAPI.Post("/unpublish/{id}", slots.UnPublishHandler())
	xmlAPI.Post("/updateUser", users.UpdateUserHandler())

	xmlAPI.Post("/showModerated", moderation.ShowModeratedHandler())
	xmlAPI.Post("/filter", moderation.FilterHandler())

	xmlAPI.Post("/uploadPhoto", photos.UploadPhoto())
	xmlAPI.Post("/photos/by", photos.GetPhotosBy())

	xmlAPI.Post("/showNotUploaded", resources.ShowNotUploadedHandler())
	xmlAPI.Post("/filterResources", resources.ShowNotUploadedHandler())

	//Stubby, mc stub face
	xmlAPI.Get("/promotions", controllers.StubEndpoint())
	xmlAPI.Get("/user/{username}/playlists", controllers.StubEndpoint())

}
