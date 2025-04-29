package ui

import (
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/components"
	. "maragu.dev/gomponents/html"
)

type PageOptions struct {
	Title            string
	ServerPathPrefix string
}

func page(o PageOptions, children ...Node) Node {
	return HTML5(HTML5Props{
		Title:    o.Title,
		Language: "en",
		Head: []Node{
			Meta(Charset("UTF-8")),
			Meta(Name("viewport"), Content("width=device-width, initial-scale=1.0, maximum-scale=1.0, viewport-fit=cover")),
			Link(Rel("icon"), Href(o.ServerPathPrefix+"/icons/favicon.ico"), Attr("sizes", "any")),
			Link(Rel("apple-touch-icon"), Href(o.ServerPathPrefix+"/icons/apple-touch-icon-180x180.png")),
			Link(Rel("stylesheet"), Href(o.ServerPathPrefix+"/css/ui-v4.1.0.css")),
			Link(Rel("stylesheet"), Href(o.ServerPathPrefix+"/css/style.css")),
			Link(Rel("manifest"), Href(o.ServerPathPrefix+"/manifest.json")),
			Script(Src(o.ServerPathPrefix + "/js/api.js")),
			Script(Src(o.ServerPathPrefix + "/js/ws.js")),
			Script(Src(o.ServerPathPrefix + "/js/utils.js")),
			Script(Raw(`window.utils.registerServiceWorker()`)),
		},
		Body: []Node{
			Group(children),
		},
	})
}
