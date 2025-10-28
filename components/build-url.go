package components

import (
	"github.com/a-h/templ"
	"github.com/knackwurstking/picow-led/env"
)

func buildURL(path string) string {
	return env.Args.ServerPathPrefix + path
}

func HXHomeSectionDevices() templ.SafeURL {
	return templ.SafeURL(buildURL("/htmx/home/section/devices"))
}

func HXHomeSectionGroups() templ.SafeURL {
	return templ.SafeURL(buildURL("/htmx/home/section/groups"))
}

func HXEditDeviceDialog() templ.SafeURL {
	return templ.SafeURL(buildURL("/htmx/dialog/edit-device"))
}

func HXEditGroupDialog() templ.SafeURL {
	return templ.SafeURL(buildURL("/htmx/dialog/edit-group"))
}
