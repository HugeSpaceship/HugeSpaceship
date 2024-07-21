package web_api

import (
	"github.com/HugeSpaceship/HugeSpaceship/internal/config"
	v3 "github.com/HugeSpaceship/HugeSpaceship/internal/http/api/web_api/v3"
	resMan "github.com/HugeSpaceship/HugeSpaceship/internal/resources"
	"github.com/go-chi/chi/v5"
)

func APIBootstrap(cfg *config.Config, res *resMan.ResourceManager) func(r chi.Router) {
	return func(r chi.Router) {
		r.Route("/v3", v3.V3APIBootstrap())
	}

}
