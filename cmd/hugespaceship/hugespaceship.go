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
	"log/slog"
	"net/http"
	"os"
)

// main is the entrypoint for the server
func main() {
	cfg, err := config.LoadConfig(false)
	if err != nil {
		slog.Error("Failed to load config", "err", err)
		os.Exit(1)
	}

	err = logger.LoggingInit(cfg)
	if err != nil {
		slog.Error("Failed to initialize logger", "err", err)
		// Arguably log formatting failing to work properly is non-fatal
	}

	pool := db.Open(cfg)            // Open a connection to the DB
	err = migration.MigrateDB(pool) // Migrate the DB to the latest schema
	if err != nil {
		slog.Error("Failed to migrate DB", "err", err)
		os.Exit(1)
	}

	res, err := resources.NewResourceManager(pool, cfg)
	if err != nil {
		slog.Error("Failed to create resource manager", "err", err)
		os.Exit(1)
	}

	// Init the router
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// everything starts at /api
	r.Route("/api/LBP_XML", func(r chi.Router) {
		// LittleBigPlanet compatible GameAPI
		game_api.APIBootstrap(r, cfg, res, pool)

		// Resource server
		game_api.ResourceBootstrap(r, res, pool)
	})

	r.Route("/api", web_api.APIBootstrap(pool))

	slog.Info("Server started", "listenAddr", cfg.ListenAddr)
	err = http.ListenAndServe(cfg.ListenAddr, r)
	if err != nil {
		panic(err)
	}
}
