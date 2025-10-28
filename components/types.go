package components

import "github.com/a-h/templ"

type ID string

type HXProps struct {
	Get               templ.SafeURL
	Target            string
	EnableTriggerLoad bool
}
