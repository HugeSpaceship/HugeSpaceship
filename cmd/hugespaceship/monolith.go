package main

/*
	The Monolith server contains all the individual services as one, this is to aid in development.
	It's also for smaller instances where scalability is not yet an issue

*/

import (
	_ "embed"
	"github.com/HugeSpaceship/HugeSpaceship/internal/config"
	"github.com/HugeSpaceship/HugeSpaceship/internal/db"
	"github.com/HugeSpaceship/HugeSpaceship/internal/db/migration"
	"github.com/HugeSpaceship/HugeSpaceship/internal/http/api/game_api"
	"github.com/HugeSpaceship/HugeSpaceship/internal/http/api/web_api"
	"github.com/HugeSpaceship/HugeSpaceship/internal/logger"
	"github.com/HugeSpaceship/HugeSpaceship/internal/resources"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"
	"log/slog"
	"net/http"
	"strconv"
)

// main is the entrypoint for the server
func main() {
	v := config.LoadConfig(false)

	logger.LoggingInit("hugespaceship", v)

	pool := db.Open(v)               // Open a connection to the DB
	err := migration.MigrateDB(pool) // Migrate the DB to the latest schema
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to migrate database")
	}

	res := resources.NewResourceManager(v, pool)
	err = res.Start()
	if err != nil {
		slog.Error("Failed to start resource manager", slog.Any("err", err))
	}

	// Init the router
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// everything starts at /api
	r.Route("/api/LBP_XML", func(r chi.Router) {
		// LittleBigPlanet compatible API
		game_api.APIBootstrap(r, v, res)

		// Resource server
		game_api.ResourceBootstrap(r, v, res)
	})

	r.Route("/api", web_api.APIBootstrap(v, res))

	slog.Info("Server started", "port", strconv.Itoa(v.GetInt("http.port")))
	err = http.ListenAndServe("0.0.0.0:"+strconv.Itoa(v.GetInt("http.port")), r)
	if err != nil {
		panic(err)
	}
}
