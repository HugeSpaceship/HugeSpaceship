package common

import "github.com/HugeSpaceship/HugeSpaceship/internal/http/website/theming"

type Info struct {
	InstanceName  string
	InstanceTheme *theming.Theme
	Debug         bool
}
