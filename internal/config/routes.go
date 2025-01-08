package config

var themeViews = []string{"web/views/theme/base.html"}

// Returns the themeViews slice.
func UiLayoutViews() []string {
	return themeViews[:]
}
