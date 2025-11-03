package components

import (
	"github.com/a-h/templ"
)

type ID string

type HXProps struct {
	URL               templ.SafeURL
	Method            string
	Target            string
	Swap              string
	Trigger           string
	EnableLoadTrigger bool
	BeforeRequest     templ.ComponentScript
	AfterRequest      templ.ComponentScript
}

func (hx *HXProps) Attributes() templ.Attributes {
	attributes := map[string]any{
		"hx-target":                 hx.Target,
		"hx-swap":                   hx.getSwap(),
		"hx-trigger":                hx.getTrigger(),
		"hx-on:htmx:before-request": hx.BeforeRequest,
		"hx-on:htmx:after-request":  hx.AfterRequest,
	}

	switch hx.Method {
	case "GET":
		attributes["hx-get"] = hx.URL
	case "POST":
		attributes["hx-post"] = hx.URL
	case "PUT":
		attributes["hx-put"] = hx.URL
	case "DELETE":
		attributes["hx-delete"] = hx.URL
	default:
		panic("unsupported method")
	}

	return attributes
}

func (hx HXProps) getSwap() string {
	if hx.Swap == "" {
		return "innerHTML"
	}
	return hx.Swap
}

func (hx *HXProps) getTrigger() string {
	if hx.Trigger == "" && !hx.EnableLoadTrigger {
		panic("no trigger events provided")
	}

	if hx.EnableLoadTrigger {
		return "load, " + hx.Trigger
	}

	return hx.Trigger
}
