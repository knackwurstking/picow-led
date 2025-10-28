package components

import (
	"strings"

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
		"hx-trigger": hx.getTrigger("click"),
	}
}

func (hx HXProps) getSwap() string {
	if hx.Swap == "" {
		return "innerHTML"
	}
	return hx.Swap
}

func (hx *HXProps) getTrigger(events ...string) string {
	if len(events) == 0 && !hx.EnableLoadTrigger {
		panic("no trigger events provided")
	}

	if hx.EnableLoadTrigger {
		return "load, " + strings.Join(events, ", ")
	}

	return strings.Join(events, ", ")
}
