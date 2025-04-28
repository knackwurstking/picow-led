package ui

import (
	"encoding/json"
	"fmt"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/components"
	. "maragu.dev/gomponents/html"
)

func page(title string, serverPathPrefix string, children ...Node) Node {
	return HTML5(HTML5Props{
		Title:    title,
		Language: "en",
		Head: []Node{
			Meta(Charset("UTF-8")),
			Meta(Name("viewport"), Content("width=device-width, initial-scale=1.0, maximum-scale=1.0, viewport-fit=cover")),
			Link(Rel("icon"), Href(serverPathPrefix+"/icons/favicon.ico"), Attr("sizes", "any")),
			Link(Rel("apple-touch-icon"), Href(serverPathPrefix+"/icons/apple-touch-icon-180x180.png")),
			Link(Rel("stylesheet"), Href(serverPathPrefix+"/css/ui-v4.1.0.css")),
			Link(Rel("stylesheet"), Href(serverPathPrefix+"/css/style.css")),
			Link(Rel("manifest"), Href(serverPathPrefix+"/manifest.json")),
			Script(Src(serverPathPrefix + "/js/api.js")),
			Script(Src(serverPathPrefix + "/js/ws.js")),
			Script(Src(serverPathPrefix + "/js/utils.js")),
			Script(Raw(`window.utils.registerServiceWorker()`)),
		},
		Body: []Node{
			Group(children),
		},
	})
}

type basePageLayoutOptions struct {
	Title            string
	AppBarTitle      string
	ServerPathPrefix string

	EnableBackButton   bool
	BackButtonCallback string

	EnableGoToSettingsButton bool
}

// basePageLayout will call page(...)
func basePageLayout(o basePageLayoutOptions, children ...Node) Node {
	return page(o.Title, o.ServerPathPrefix,
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

					onlineIndicator(false),
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
						Attr("data-ui-color", "secondary"),
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
		),
	)
}

func toJSON(data any) []byte {
	d, _ := json.Marshal(data)
	return d
}
