package ui

import (
	"fmt"
	"picow-led/internal/api"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func DevicesAddrPage(serverPathPrefix string, d *api.Device) Node {
	appBarTitle := d.Server.Addr
	if d.Server.Name != "" {
		appBarTitle = d.Server.Name
	}

	return basePageLayout(
		LayoutBaseOptions{
			PageOptions: PageOptions{
				Title:            fmt.Sprintf("PicoW LED | %s", d.Server.Addr),
				ServerPathPrefix: serverPathPrefix,
			},
			AppBarTitle:        appBarTitle,
			EnableBackButton:   true,
			BackButtonCallback: fmt.Sprintf("location.pathname = \"%s\"", serverPathPrefix),
		},

		Div(
			Class("ui-container ui-auto-scroll ui-hide-scrollbar"),
			Style("height: 100%; padding-top: var(--ui-app-bar-height);"),

			// TODO: Color cache

			// TODO: Color picker
		),
	)
}
