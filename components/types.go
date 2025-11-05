package components

import (
	"net/http"

	"github.com/a-h/templ"
)

type HXProps struct {
	URL           templ.SafeURL
	Method        string
	Target        string
	Swap          string
	Trigger       string
	BeforeRequest string
	AfterRequest  string
}

func (hx *HXProps) Attributes() templ.Attributes {
	attributes := map[string]any{
		"hx-target":  hx.Target,
		"hx-swap":    hx.getSwap(),
		"hx-trigger": hx.getTrigger(),
	}

	if hx.BeforeRequest != "" {
		attributes["hx-on:htmx:before-request"] = hx.BeforeRequest
	}
	if hx.AfterRequest != "" {
		attributes["hx-on:htmx:after-request"] = hx.AfterRequest
	}

	switch hx.Method {
	case http.MethodGet:
		attributes["hx-get"] = string(hx.URL)
	case http.MethodPost:
		attributes["hx-post"] = string(hx.URL)
	case http.MethodPut:
		attributes["hx-put"] = string(hx.URL)
	case http.MethodDelete:
		attributes["hx-delete"] = string(hx.URL)
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
	if hx.Trigger == "" {
		panic("no trigger events provided")
	}

	return hx.Trigger
}
