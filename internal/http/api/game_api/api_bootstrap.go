package game_api

import (
	"github.com/HugeSpaceship/HugeSpaceship/internal/config"
	"github.com/HugeSpaceship/HugeSpaceship/internal/http/api/game_api/controllers"
	"github.com/HugeSpaceship/HugeSpaceship/internal/http/api/game_api/controllers/auth"
	"github.com/HugeSpaceship/HugeSpaceship/internal/http/api/game_api/controllers/match"
	"github.com/HugeSpaceship/HugeSpaceship/internal/http/api/game_api/controllers/moderation"
	"github.com/HugeSpaceship/HugeSpaceship/internal/http/api/game_api/controllers/photos"
	"github.com/HugeSpaceship/HugeSpaceship/internal/http/api/game_api/controllers/resources"
	"github.com/HugeSpaceship/HugeSpaceship/internal/http/api/game_api/controllers/settings"
	"github.com/HugeSpaceship/HugeSpaceship/internal/http/api/game_api/controllers/slots"
	"github.com/HugeSpaceship/HugeSpaceship/internal/http/api/game_api/controllers/users"
	"github.com/HugeSpaceship/HugeSpaceship/internal/http/api/game_api/middlewares"
	"github.com/HugeSpaceship/HugeSpaceship/internal/http/middleware"
	resMan "github.com/HugeSpaceship/HugeSpaceship/internal/resources"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/jackc/pgx/v5/pgxpool"
)

func ResourceBootstrap(group chi.Router, res *resMan.ResourceManager, pool *pgxpool.Pool) {

	group.With(middleware.DBCtxMiddleware(pool), middlewares.TicketAuthMiddleware).Get("/r/{hash}", resources.GetResourceHandler(res))
}

func APIBootstrap(r chi.Router, cfg *config.Config, res *resMan.ResourceManager, pool *pgxpool.Pool) {
	r.Use(middlewares.PSPVersionMiddleware, middlewares.ServerHeaderMiddleware)

	r.Use(chiMiddleware.RequestID)
	r.Use(chiMiddleware.Recoverer)
	r.Use(chiMiddleware.NoCache)

	r.With(middleware.DBCtxMiddleware(pool)).Post("/login", auth.LoginHandler())
	r.Get("/eula", controllers.EulaHandler())

	// LittleBigPlanet compatible GameAPI, required NpTicket auth
	authGameAPI := r.With(middleware.DBCtxMiddleware(pool), middlewares.TicketAuthMiddleware)
	authGameAPI.With(middleware.DBCtxMiddleware(pool)).Get("/announce", controllers.AnnounceHandler())
	authGameAPI.Get("/network_settings.nws", settings.NetSettingsHandler())

	// LittleBigPlanet compatible GameAPI with digest calculation
	digestRequiredAPI := authGameAPI.With(middlewares.DigestMiddleware(cfg))
	digestRequiredAPI.Post("/match", match.MatchEndpoint())
	digestRequiredAPI.Post("/upload/{hash}", resources.UploadResources(res))

	// GameAPI with automatic content type
	xmlAPI := digestRequiredAPI.With(render.SetContentType(render.ContentTypeXML))
	xmlAPI.Get("/user/{username}", users.UserGetHandler())

	xmlAPI.Post("/npdata", settings.NpDataEndpoint())
	xmlAPI.Get("/notification", controllers.NotificationController()) // Stub
	xmlAPI.Post("/goodbye", auth.LogoutHandler())
	xmlAPI.Get("/news", controllers.NewsHandler())
	xmlAPI.Get("/news/{id}", controllers.LBP2NewsHandler())
	xmlAPI.Get("/stream", controllers.StreamHandler())

	xmlAPI.Get("/slots", slots.GetSlotsHandler())
	xmlAPI.Get("/slots/by", slots.GetSlotsByHandler())
	xmlAPI.Get("/slots/lbp2luckydip", slots.LuckyDipHandler())
	xmlAPI.Get("/slots/thumbs", slots.HighestRatedLevelsHandler())
	xmlAPI.Get("/slots/lbp2cool", slots.HighestRatedLevelsHandler())
	xmlAPI.Get("/slots/cool", slots.HighestRatedLevelsHandler())
	xmlAPI.Get("/slots/highestRated", slots.HighestRatedLevelsHandler())

	xmlAPI.Get("/s/user/{id}", slots.GetSlotHandler())

	xmlAPI.Post("/scoreboard/{levelType}/{levelID}", slots.UploadScoreHandler())

	xmlAPI.Post("/startPublish", slots.StartPublishHandler(res))
	xmlAPI.Post("/publish", slots.PublishHandler())
	xmlAPI.Post("/unpublish/{id}", slots.UnPublishHandler())
	xmlAPI.Post("/updateUser", users.UpdateUserHandler())

	xmlAPI.Post("/showModerated", moderation.ShowModeratedHandler())
	xmlAPI.Post("/filter", moderation.FilterHandler())

	xmlAPI.Post("/uploadPhoto", photos.UploadPhoto())
	xmlAPI.Post("/photos/by", photos.GetPhotosBy())

	xmlAPI.Post("/showNotUploaded", resources.ShowNotUploadedHandler(res))
	xmlAPI.Post("/filterResources", resources.ShowNotUploadedHandler(res))

	// FIXME: Stubbed APIs
	// These need to be implemented at some point
	xmlAPI.Get("/promotions", controllers.StubEndpoint())
	xmlAPI.Get("/user/{username}/playlists", controllers.StubEndpoint())

}
