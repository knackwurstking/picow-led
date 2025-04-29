package ui

import (
	"fmt"
	"net/url"
	"picow-led/internal/api"
	"picow-led/ui/components"
	"picow-led/ui/utils"
	"slices"
	"strconv"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func DevicesPage(serverPathPrefix string, devices ...*api.Device) Node {
	deviceListItemOptions := []deviceListItemOption{}
	for _, d := range devices {
		deviceListItemOptions = append(deviceListItemOptions, deviceListItemOption{
			Device:           d,
			serverPathPrefix: serverPathPrefix,
		})
	}

	return basePageLayout(
		LayoutBaseOptions{
			PageOptions: PageOptions{
				Title:            "PicoW LED | Devices",
				ServerPathPrefix: serverPathPrefix,
			},
			AppBarTitle:              "Devices",
			EnableGoToSettingsButton: true,
		},

		Div(
			Class("ui-container ui-auto-scroll ui-hide-scrollbar"),
			Style("height: 100%; padding-top: var(--ui-app-bar-height);"),

			// Devices List
			Span(
				Class("ui-flex column gap align-center"),
				Map(deviceListItemOptions, deviceListItem),
			),
		),
	)
}

type deviceListItemOption struct {
	Device           *api.Device
	serverPathPrefix string
}

func deviceListItem(o deviceListItemOption) Node {
	colorS := []string{}
	powerButtonState := "off"
	if o.Device.Color != nil {
		for _, c := range o.Device.Color[:3] {
			colorS = append(colorS, strconv.Itoa(int(c)))
		}

		if slices.Max(o.Device.Color) > 0 {
			powerButtonState = "on"
		}
	}

	var name string
	if o.Device.Server.Name != "" {
		name = o.Device.Server.Name
	} else {
		name = o.Device.Server.Addr
	}

	return Section(
		Class("device-list-item ui-flex row gap justify-between align-center ui-padding"),
		Style("width: 100%;"),
		Attr("data-ui-theme", "dark"),
		Attr("data-json", string(utils.ToJSON(o.Device))),

		H3(
			Class("title ui-padding"),
			Text(name),
		),

		Span(
			Class("ui-flex-item ui-flex row gap"),
			Style("flex: 0;"),

			Span(
				Class("ui-flex-item"),
				Style("flex: 0;"),

				Button(
					Attr("data-ui-variant", "ghost"),
					Attr("data-ui-color", "secondary"),
					Attr("onclick", fmt.Sprintf(
						"location.pathname = \"%s/devices/%s\"",
						o.serverPathPrefix,
						url.QueryEscape(o.Device.Server.Addr),
					)),

					Text("Edit"),
				),
			),

			Span(
				Class("ui-flex-item"),
				Style("flex: 0;"),

				components.PowerButton(powerButtonState, colorS),
			),
		),
	)
}
