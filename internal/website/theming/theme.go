package theming

import "html/template"

// Theme defines a server-owner editable theme for the website
type Theme struct {
	// Version is used to check if a theme is outdated and will likely not work with the current server version
	Version uint16 `json:"version"`
	// The id of the theme, like com.example.theme-name
	ID string `json:"id"`
	// ThemeVersion
	ThemeVersion string `json:"theme-version"`
	// The displayed name of the theme
	Name string `json:"name"`
	// Some info about the theme
	Description string `json:"description"`
	// Who made the theme
	Author string `json:"author"`
	// Path to the theme
	Path string `json:"-"`
	// List of stylesheets to be loaded by the client, paths are relative to the css folder in the theme
	Stylesheets []string `json:"stylesheets"`
	// If a theme is built-in then it is loaded as an embedded resource, instead of as a file
	BuiltIn bool `json:"-"`
	// If the theme replaces the built-in templates
	OverrideTemplates bool `json:"override-templates"`
	// The template associated with the theme, this is used so that the theme can override parts of the default template.
	Template *template.Template `json:"-"`
}

// BuiltInThemes is a list of the themes that come pre-installed with HugeSpaceship
var BuiltInThemes = []Theme{
	{
		ID:          "builtin.hugespaceship.shuttle",
		Name:        "Shuttle",
		Description: "The default HugeSpaceship theme",
		Author:      "The HugeSpaceship contributors",
		Path:        "shuttle",
		Stylesheets: []string{"theme.css"},
		BuiltIn:     true,
	},
	{
		ID:          "builtin.hugespaceship.lighthouse",
		Name:        "Lighthouse",
		Description: "A ProjectLighthouse inspired theme",
		Author:      "The HugeSpaceship contributors, uses assets from ProjectLighthouse",
		Path:        "lighthouse",
		Stylesheets: []string{"theme.css"},
		BuiltIn:     true,
	},
}
