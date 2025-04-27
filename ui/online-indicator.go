package ui

import (
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

// onlineIndicator styles in "public/css/style.css"
func onlineIndicator(state bool) Node {
	return Span(
		Class("online-indicator"),
		Style("--mono: 1;"),
		If(state, Attr("data-state", "online")),
		If(!state, Attr("data-state", "offline")),
	)
}
