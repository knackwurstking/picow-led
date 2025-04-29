package ui

import (
	"fmt"
	"picow-led/ui/components"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

type LayoutBaseOptions struct {
	PageOptions

	AppBarTitle string

	EnableBackButton   bool
	BackButtonCallback string

	EnableGoToSettingsButton bool
}

// basePageLayout will call page(...)
func basePageLayout(o LayoutBaseOptions, children ...Node) Node {
	return page(
		PageOptions{
			Title:            o.Title,
			ServerPathPrefix: o.ServerPathPrefix,
		},

		Main(
			Group(children),

			// UI App Bar
			Div(
				Class("ui-app-bar"),
				Attr("data-ui-position", "top"),

				Span(
					Class("ui-app-bar-left"),

					If(o.EnableBackButton, Button(
						Attr("data-ui-variant", "ghost"),

						If(
							o.BackButtonCallback != "",
							Attr("onclick", o.BackButtonCallback),
						),

						Text("Back"),
					)),

					components.OnlineIndicator(false),
				),

				Span(
					Class("ui-app-bar-center"),

					If(o.AppBarTitle != "",
						Span(
							H3(Text(o.AppBarTitle)),
						),
					),
				),

				Span(
					Class("ui-app-bar-right"),

					If(o.EnableGoToSettingsButton, Button(
						Attr("data-ui-variant", "ghost"),
						Attr("onclick", fmt.Sprintf(
							"location.pathname = \"%s/settings\"",
							o.ServerPathPrefix,
						)),

						Text("Settings"),
					)),
				),
			),

			// Scripts section
			Script(
				Raw("window.utils.setOnlineIndicator(true)"),
			),
			Script(Src(o.ServerPathPrefix+"/js/base-layout.js")),
		),
	)
}
