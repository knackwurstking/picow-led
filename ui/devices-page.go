package ui

import (
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func DevicesPage(serverPathPrefix string) Node {
	return page("PicoW LED | Devices", serverPathPrefix,
		Main(
			// UI App Bar
			Div(
				Class("ui-app-bar"),
				Span(
					Class("ui-app-bar-left"),
					// onlineIndicator(),
				),
				Span(
					Class("ui-app-bar-center"),
				),
				Span(
					Class("ui-app-bar-right"),
				),
			),

			// Devices List
			Span(
				Class("ui-flex column gap align-center"),
				// TODO: Iter devices
			),

			// Templates for later
			// deviceListItemTemplate()
		),
	)
}
