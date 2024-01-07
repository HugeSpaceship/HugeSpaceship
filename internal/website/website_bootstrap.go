package website

import (
	"HugeSpaceship/internal/config"
	"HugeSpaceship/internal/model/common"
	"HugeSpaceship/internal/website/pages"
	"HugeSpaceship/internal/website/theming"
	"embed"
	"github.com/gin-gonic/gin"
	"io/fs"
	"net/http"
)

//go:embed static
var staticFiles embed.FS

func Bootstrap(ctx *gin.Engine, cfg *config.Config) {
	info := common.Info{
		InstanceName: "HugeSpaceship DEV",
		Debug:        cfg.Log.Debug,
	}

	if cfg.Website.UseEmbeddedResources {
		static, err := fs.Sub(staticFiles, "static")
		if err != nil {
			panic(err)
		}
		ctx.StaticFS("/static", http.FS(static))
	} else {
		ctx.Static("/static", cfg.Website.WebRoot)
	}

	themes, err := theming.LoadThemes(cfg.Website.ThemePath, ctx)
	if err != nil {
		panic(err) // Should only ever be a dir issue, should probably do some error handling here though
	}

	var exists bool
	info.InstanceTheme, exists = themes.GetTheme(cfg.Website.DefaultTheme)
	if !exists {
		info.InstanceTheme, _ = themes.GetTheme("builtin.hugespaceship.shuttle")
	}

	ctx.GET("/", pages.HomePage(info))
	ctx.GET("/earth", pages.EarthPage(info))
}
