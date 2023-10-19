package theming

func GetTheme(name string) *Theme {

	for _, theme := range BuiltInThemes {
		if theme.ID == name {
			return &theme
		}
	}
	return nil
}
