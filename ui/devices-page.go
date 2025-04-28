package ui

import (
	"picow-led/internal/api"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func DevicesPage(serverPathPrefix string, devices ...*api.Device) Node {
	return page("PicoW LED | Devices", serverPathPrefix,
		Main(
			// UI App Bar
			Div(
				Class("ui-app-bar"),
				Span(
					Class("ui-app-bar-left"),
					onlineIndicator(false),
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
				Map(devices, deviceListItem),
			),

			// Templates for later
			// deviceListItemTemplate()

			// Scripts section
			Script(
				Raw("window.utils.setOnlineIndicator(true)"),
			),
		),
	)
}

func deviceListItem(d *api.Device) Node {
	return Section(
		Class("device-list-item ui-flex row gap justify-between align-center ui-padding"),
		Style("width: 100%;"),
		If(d.Server.Name != "", H4(Text(d.Server.Name))),
		If(d.Server.Name == "", H4(Text(d.Server.Addr))),
		Span(
			Class("ui-flex-item"),
			Style("flex: 0;"),
			powerButton(d),
		),
	)
}
