package components

import (
	"github.com/a-h/templ"
)

type ID string

type HXProps struct {
	Get               templ.SafeURL
	Target            string
	Swap              string
	Trigger           string
	EnableLoadTrigger bool
}

func (hx *HXProps) Attributes() templ.Attributes {
	return map[string]any{
		"hx-get":     string(hx.Get),
		"hx-target":  hx.Target,
		"hx-swap":    hx.getSwap(),
		"hx-trigger": hx.getTrigger(),
	}
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
