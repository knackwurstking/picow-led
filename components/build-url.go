package components

import (
	"github.com/a-h/templ"
	"github.com/knackwurstking/picow-led/env"
)

func BuildURL(path string) string {
	return env.Args.ServerPathPrefix + path
}

func HXHomeSectionDevices() templ.SafeURL {
	return templ.SafeURL(BuildURL("/htmx/home/section/devices"))
}

func HXHomeSectionGroups() templ.SafeURL {
	return templ.SafeURL(BuildURL("/htmx/home/section/groups"))
}

func HXEditDeviceDialog() templ.SafeURL {
	return templ.SafeURL(BuildURL("/htmx/dialog/edit-device"))
}

func HXEditGroupDialog() templ.SafeURL {
	return templ.SafeURL(BuildURL("/htmx/dialog/edit-group"))
}
