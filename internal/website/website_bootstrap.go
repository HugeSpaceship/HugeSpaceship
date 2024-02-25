package website

import (
	"HugeSpaceship/internal/config"
	"HugeSpaceship/internal/model/common"
	"HugeSpaceship/internal/website/api"
	"HugeSpaceship/internal/website/pages"
	"HugeSpaceship/internal/website/theming"
	"embed"
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
		r.Get("/", pages.HomePage(info))
		r.Get("/", pages.HomePage(info))

		// API
		r.Get("/api/v1/test", api.SlotAPI(info))
		r.Get("/api/v1/image/{hash}", api.ImageConverterHandler(cfg))
	}
}
