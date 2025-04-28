package ui

import (
	"picow-led/internal/api"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func SettingsPage(serverPathPrefix string, devices ...*api.Device) Node {
	return basePageLayout("PicoW LED | Settings", serverPathPrefix,
		Main(
		// TODO: Settings page content here
		),
	)
}
