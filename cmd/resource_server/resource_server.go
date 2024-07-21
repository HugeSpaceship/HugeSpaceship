package main

import (
	"github.com/HugeSpaceship/HugeSpaceship/internal/config"
	"github.com/HugeSpaceship/HugeSpaceship/internal/db"
	"github.com/HugeSpaceship/HugeSpaceship/internal/db/migration"
	"github.com/HugeSpaceship/HugeSpaceship/internal/http/api/game_api"
	"github.com/HugeSpaceship/HugeSpaceship/internal/logger"
	"github.com/HugeSpaceship/HugeSpaceship/internal/resources"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
)

func main() {
	v := config.LoadConfig(false)

	logger.LoggingInit("api_server", v)

	pool := db.Open(v)               // Open a connection to the DB
	err := migration.MigrateDB(pool) // Migrate the DB to the latest schema
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to migrate database")
	}

	res := resources.NewResourceManager(v)

	// Initialize chi router
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// everything starts at /api
	r.Route("/api/LBP_XML", func(r chi.Router) {
		// LittleBigPlanet compatible API
		game_api.ResourceBootstrap(r, v, res)
	})

	err = http.ListenAndServe("0.0.0.0:"+strconv.Itoa(v.GetInt("http.port")), r)
	if err != nil {
		panic(err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
}
