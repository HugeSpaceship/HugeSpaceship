package website

import (
	"HugeSpaceship/internal/config"
	"HugeSpaceship/internal/model/common"
	"HugeSpaceship/internal/website/pages"
	"HugeSpaceship/internal/website/theming"
	"embed"
	"github.com/gin-gonic/gin"
)

//go:embed static
var staticFiles embed.FS

func Bootstrap(ctx *gin.Engine, cfg *config.Config) {
	// TODO: Embed files by default, but allow loading files locally if the user wants to, will be useful for debugging

	theme := theming.GetTheme("builtin.hugespaceship.shuttle")
	if theme == nil {
		panic("Invalid theme")
	}

	info := common.Info{
		InstanceName:  "HugeSpaceship DEV",
		InstanceTheme: theme,
	}

	ctx.LoadHTMLGlob("./internal/website/partials/*")
	ctx.Static("/static", "./internal/website/static")
	ctx.GET("/", pages.HomePage(info))
	ctx.GET("/earth", pages.EarthPage(info))
}
