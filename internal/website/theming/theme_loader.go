package theming

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"os"
	"path"
)

// ThemeVersion defines the current version of the theme schema,
// if changes are made to the way that themes are structured please increment this by 1.
const ThemeVersion = 2

// GetTheme retrieves the theme with the specified ID, and a boolean indicating whether the theme exists
func (t Themes) GetTheme(id string) (*Theme, bool) {
	for _, theme := range BuiltInThemes {
		if theme.ID == id {
			return &theme, true
		}
	}

	for _, theme := range t {
		if theme.ID == id {
			return theme, true
		}
	}

	return nil, false
}

type Themes []*Theme

func LoadThemes(themesPath string, ctx *gin.Engine) (themes Themes, err error) {
	themesDir, err := os.ReadDir(themesPath)

	for _, themeDir := range themesDir {
		if !themeDir.IsDir() { // TODO: add support for themesDir in zip files
			continue
		}

		themeJson, err := os.ReadFile(path.Clean(themesPath) + "/" + themeDir.Name() + "/theme.json")
		if err != nil {
			log.Warn().Str("theme-dir", themeDir.Name()).Msg("Could not read theme.json, is this a valid theme?")
			continue
		}

		theme := new(Theme)
		err = json.Unmarshal(themeJson, theme)
		if err != nil {
			log.Warn().Str("theme-dir", themeDir.Name()).Msg("Could not parse theme.json, is this a valid theme?")
			continue
		}

		ctx.Static("/theming/"+theme.ID, path.Clean(themesPath)+"/"+themeDir.Name())
		theme.Path = themeDir.Name()
		themeTempl, err := LoadTemplateOverrides(theme, themesPath)
		if err == nil {
			theme.Template = themeTempl
		} else {
			log.Warn().Str("theme", theme.ID).Err(err).Msg("Failed to load theme templates, using defaults")
			theme.Template, err = LoadBaseTemplates()
			if err != nil {
				panic(err)
			}
		}

		themes = append(themes, theme)
	}
	return
}
