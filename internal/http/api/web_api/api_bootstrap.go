package web_api

import (
	v3 "github.com/HugeSpaceship/HugeSpaceship/internal/http/api/web_api/v3"
	resMan "github.com/HugeSpaceship/HugeSpaceship/internal/resources"
	"github.com/go-chi/chi/v5"
	"github.com/spf13/viper"
)

func APIBootstrap(v *viper.Viper, res *resMan.ResourceManager) func(r chi.Router) {
	return func(r chi.Router) {
		r.Route("/v3", v3.V3APIBootstrap())
	}

}
