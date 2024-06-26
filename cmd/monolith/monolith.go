package main

/*
	The Monolith server contains all the individual services as one, this is to aid in development.
	It's also for smaller instances where scalability is not yet an issue

*/

import (
	"HugeSpaceship/internal/config"
	"HugeSpaceship/internal/http/api/game_api"
	"HugeSpaceship/internal/http/website"
	"HugeSpaceship/pkg/db"
	"HugeSpaceship/pkg/db/migration"
	"HugeSpaceship/pkg/logger"
	_ "embed"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"
	"net/http"
	"strconv"
)

// main is the entrypoint for the server
func main() {
	cfg, err := config.LoadConfig(false)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	logger.LoggingInit("hugespaceship", cfg)

	pool := db.Open(cfg)            // Open a connection to the DB
	err = migration.MigrateDB(pool) // Migrate the DB to the latest schema
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to migrate database")
	}

	// Init the router
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// everything starts at /api
	r.Route("/api/LBP_XML", func(r chi.Router) {
		// LittleBigPlanet compatible API
		game_api.APIBootstrap(r, cfg)

		// Resource server
		game_api.ResourceBootstrap(r, cfg)
	})

	r.Group(website.Bootstrap(cfg))

	err = http.ListenAndServe("0.0.0.0:"+strconv.Itoa(cfg.HTTPPort), r)
	if err != nil {
		panic(err)
	}
}
