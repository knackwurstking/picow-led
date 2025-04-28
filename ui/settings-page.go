package ui

import (
	"fmt"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func SettingsPage(serverPathPrefix string) Node {
	return basePageLayout(
		basePageLayoutOptions{
			Title:              "PicoW LED | Settings",
			AppBarTitle:        "Settings",
			ServerPathPrefix:   serverPathPrefix,
			EnableBackButton:   true,
			BackButtonCallback: fmt.Sprintf("location.pathname = \"%s\"", serverPathPrefix),
		},

		Main(
			Class("ui-container"),
			// TODO: Settings page content here
		),
	)
}
