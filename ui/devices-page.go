package ui

import (
	"picow-led/internal/api"
	"slices"
	"strconv"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func DevicesPage(serverPathPrefix string, devices ...*api.Device) Node {
	return basePageLayout("PicoW LED | Devices", serverPathPrefix,
		Div(
			Class("ui-container ui-auto-scroll ui-hide-scrollbar"),
			Style("max-height: 100%; padding-top: var(--ui-app-bar-height);"),
			// Devices List
			Span(
				Class("ui-flex column gap align-center"),
				Map(devices, deviceListItem),
			),
		),
	)
}

func deviceListItem(d *api.Device) Node {
	colorS := []string{}
	powerButtonState := "off"
	if d.Color != nil {
		for _, c := range d.Color[:3] {
			colorS = append(colorS, strconv.Itoa(int(c)))
		}

		if slices.Max(d.Color) > 0 {
			powerButtonState = "on"
		}
	}

	var name string
	if d.Server.Name != "" {
		name = d.Server.Name
	} else {
		name = d.Server.Addr
	}

	return Section(
		Class("device-list-item ui-flex row gap justify-between align-center ui-padding"),
		Style("width: 100%;"),
		Attr("data-ui-theme", "dark"),
		Attr("data-json", string(toJSON(d))),
		H4(
			Class("title ui-outline-text ui-padding"),
			Text(name),
		),
		Span(
			Class("ui-flex-item"),
			Style("flex: 0;"),
			powerButton(powerButtonState, colorS),
		),
	)
}
