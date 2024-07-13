package website

import (
	"embed"
	"github.com/HugeSpaceship/HugeSpaceship/internal/config"
	"github.com/HugeSpaceship/HugeSpaceship/internal/http/middleware"
	"github.com/HugeSpaceship/HugeSpaceship/internal/http/website/api"
	"github.com/HugeSpaceship/HugeSpaceship/internal/http/website/pages"
	"github.com/HugeSpaceship/HugeSpaceship/internal/http/website/theming"
	"github.com/HugeSpaceship/HugeSpaceship/internal/model/common"
	"github.com/go-chi/chi/v5"
	"io/fs"
	"net/http"
)

//go:embed static
var staticFiles embed.FS

func Bootstrap(cfg *config.Config) func(r chi.Router) {
	return func(r chi.Router) {
		info := common.Info{
			InstanceName: "HugeSpaceship DEV",
			Debug:        cfg.Log.Debug,
		}

		if cfg.Website.UseEmbeddedResources {
			static, err := fs.Sub(staticFiles, "static")
			if err != nil {
				panic(err)
			}
			r.Handle("/static", http.FileServerFS(static))
		} else {
			r.Handle("/static", http.FileServer(http.Dir(cfg.Website.WebRoot)))
		}

		themes, err := theming.LoadThemes(cfg.Website.ThemePath, r)
		if err != nil {
			panic(err) // Should only ever be a dir issue, should probably do some error handling here though
		}

		var exists bool
		info.InstanceTheme, exists = themes.GetTheme(cfg.Website.DefaultTheme)
		if !exists {
			info.InstanceTheme, _ = themes.GetTheme("builtin.hugespaceship.shuttle")
		}

		// Pages
		r.With(middleware.DBCtxMiddleware).Get("/", pages.HomePage(info))

		// API
		r.Route("/api/v1", func(r chi.Router) {
			r.Use(middleware.DBCtxMiddleware)

			r.Get("/api/v1/test", api.SlotAPI(info))

			r.Route("/image/{hash}", func(r chi.Router) {
				r.Get("/", api.ImageConverterHandler(cfg))
			})

		})

	}
}
