package ui

import (
	"picow-led/internal/api"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func DevicesPage(serverPathPrefix string, devices ...*api.Device) Node {
	return basePageLayout("PicoW LED | Devices", serverPathPrefix,
		Main(
			Class("ui-container"),
			// Devices List
			Span(
				Class("ui-flex column gap align-center"),
				Map(devices, deviceListItem),
			),
		),
	)
}

// TODO: Need some highlighting for color, should work with light and dark theme
func deviceListItem(d *api.Device) Node {
	return Section(
		Class("device-list-item ui-flex row gap justify-between align-center ui-padding"),
		Style("width: 100%;"),
		Attr("data-json", string(toJSON(d))),
		If(d.Server.Name != "", H4(Text(d.Server.Name))),
		If(d.Server.Name == "", H4(Text(d.Server.Addr))),
		Span(
			Class("ui-flex-item"),
			Style("flex: 0;"),
			powerButton(),
		),
	)
}
