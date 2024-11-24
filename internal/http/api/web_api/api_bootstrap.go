package web_api

import (
	v3 "github.com/HugeSpaceship/HugeSpaceship/internal/http/api/web_api/v3"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func APIBootstrap(pool *pgxpool.Pool) func(r chi.Router) {
	return func(r chi.Router) {
		r.Route("/v3", v3.APIBootstrap(pool))
	}

}
