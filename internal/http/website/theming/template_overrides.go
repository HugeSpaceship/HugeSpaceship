package theming

import (
	"HugeSpaceship/internal/http/website/partials"
	"html/template"
	"path"
)

func LoadBaseTemplates() (*template.Template, error) {
	return template.ParseFS(partials.Partials, "html/*")
}

func LoadTemplateOverrides(theme *Theme, themePath string) (*template.Template, error) {
	if !theme.OverrideTemplates {
		return LoadBaseTemplates()
	}
	templates, err := LoadBaseTemplates()
	if err != nil {
		return nil, err
	}
	return templates.ParseGlob(path.Clean(themePath) + "/" + theme.Path + "/templates/*")
}
