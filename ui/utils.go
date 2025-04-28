package ui

import (
	"encoding/json"

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
			Link(Rel("icon"), Href(serverPathPrefix+"/icons/favicon.ico")),
			Link(Rel("icon"), Href(serverPathPrefix+"/icons/apple-touch-icon-180x180.png")),
			Link(Rel("stylesheet"), Href(serverPathPrefix+"/css/ui-v4.1.0.css")),
			Link(Rel("stylesheet"), Href(serverPathPrefix+"/css/style.css")),
			Script(Src(serverPathPrefix + "/js/api.js")),
			Script(Src(serverPathPrefix + "/js/ws.js")),
			Script(Src(serverPathPrefix + "/js/utils.js")),
		},
		Body: []Node{
			Class("ui-container"),
			Group(children),
		},
	})
}

// basePageLayout will call page(...)
func basePageLayout(title string, serverPathPrefix string, children ...Node) Node {
	return page(title, serverPathPrefix,
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

			Group(children),

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
